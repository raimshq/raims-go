// Command check-api compares the exported SDK API with an independent baseline.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
)

type listedPackage struct {
	ImportPath string
	Dir        string
	GoFiles    []string
	Export     string
}

type field struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Tag      string `json:"tag,omitempty"`
	Embedded bool   `json:"embedded,omitempty"`
}

type symbol struct {
	Kind           string   `json:"kind"`
	Type           string   `json:"type,omitempty"`
	Value          string   `json:"value,omitempty"`
	Receiver       string   `json:"receiver,omitempty"`
	TypeParameters []string `json:"type_parameters,omitempty"`
	Fields         []field  `json:"fields,omitempty"`
}

type manifest struct {
	SchemaVersion int               `json:"schema_version"`
	Module        string            `json:"module"`
	Package       string            `json:"package"`
	Symbols       map[string]symbol `json:"symbols"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	write := flag.Bool("write", false, "write the baseline after reviewing API changes")
	check := flag.Bool("check", false, "check the baseline without writing (default)")
	source := flag.String("source", ".", "SDK source directory")
	baseline := flag.String("manifest", "", "independent API baseline; defaults to contracts/go-api.json")
	flag.Parse()
	if flag.NArg() != 0 || (*write && *check) {
		return fmt.Errorf("usage: check-api [--check|--write] [--source directory] [--manifest file]")
	}
	root, err := filepath.Abs(*source)
	if err != nil {
		return err
	}
	path := *baseline
	if path == "" {
		path = filepath.Join(root, "contracts", "go-api.json")
	}
	current, err := snapshot(root)
	if err != nil {
		return err
	}
	if *write {
		data, err := json.MarshalIndent(current, "", "  ")
		if err != nil {
			return err
		}
		if err := os.WriteFile(path, append(data, '\n'), 0644); err != nil {
			return err
		}
		fmt.Printf("Wrote %d exported API symbols to %s. Review the baseline diff before distribution.\n", len(current.Symbols), path)
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read independent API baseline: %w", err)
	}
	var approved manifest
	if err := json.Unmarshal(data, &approved); err != nil {
		return fmt.Errorf("invalid API baseline: %w", err)
	}
	if approved.SchemaVersion != 2 || approved.Module != current.Module || approved.Package != current.Package {
		return fmt.Errorf("API baseline metadata differs from the SDK")
	}
	var differences []string
	for name, old := range approved.Symbols {
		value, ok := current.Symbols[name]
		if !ok {
			differences = append(differences, "removed: "+name)
		} else if !reflect.DeepEqual(old, value) {
			differences = append(differences, "changed: "+name)
		}
	}
	for name := range current.Symbols {
		if _, ok := approved.Symbols[name]; !ok {
			differences = append(differences, "added: "+name)
		}
	}
	if len(differences) > 0 {
		sort.Strings(differences)
		for _, difference := range differences {
			fmt.Fprintln(os.Stderr, difference)
		}
		return fmt.Errorf("exported API differs from the independent baseline; review changes before explicit --write")
	}
	fmt.Printf("Verified %d exported API symbols against the independent baseline.\n", len(current.Symbols))
	return nil
}

func snapshot(root string) (manifest, error) {
	command := exec.Command("go", "list", "-export", "-deps", "-json", ".")
	command.Dir = root
	data, err := command.Output()
	if err != nil {
		if failure, ok := err.(*exec.ExitError); ok {
			return manifest{}, fmt.Errorf("go list failed: %s", failure.Stderr)
		}
		return manifest{}, err
	}
	exports := map[string]string{}
	var target listedPackage
	decoder := json.NewDecoder(bytes.NewReader(data))
	for {
		var pkg listedPackage
		if err := decoder.Decode(&pkg); err != nil {
			if err == io.EOF {
				break
			}
			return manifest{}, err
		}
		if pkg.Export != "" {
			exports[pkg.ImportPath] = pkg.Export
		}
		if filepath.Clean(pkg.Dir) == filepath.Clean(root) {
			target = pkg
		}
	}
	if target.ImportPath != "github.com/raimshq/raims-go" {
		return manifest{}, fmt.Errorf("expected the Raims Go SDK module")
	}
	set := token.NewFileSet()
	var files []*ast.File
	for _, name := range target.GoFiles {
		file, err := parser.ParseFile(set, filepath.Join(root, name), nil, 0)
		if err != nil {
			return manifest{}, err
		}
		files = append(files, file)
	}
	lookup := func(path string) (io.ReadCloser, error) { return os.Open(exports[path]) }
	config := types.Config{Importer: importer.ForCompiler(set, "gc", lookup)}
	pkg, err := config.Check(target.ImportPath, set, files, nil)
	if err != nil {
		return manifest{}, err
	}
	qualifier := func(other *types.Package) string {
		if other == pkg {
			return ""
		}
		return other.Path()
	}
	typeName := func(typ types.Type) string { return types.TypeString(typ, qualifier) }
	result := manifest{SchemaVersion: 2, Module: pkg.Path(), Package: pkg.Name(), Symbols: map[string]symbol{}}
	for _, name := range pkg.Scope().Names() {
		if !ast.IsExported(name) {
			continue
		}
		object := pkg.Scope().Lookup(name)
		switch value := object.(type) {
		case *types.Const:
			result.Symbols[name] = symbol{Kind: "const", Type: typeName(value.Type()), Value: value.Val().ExactString()}
		case *types.Var:
			result.Symbols[name] = symbol{Kind: "var", Type: typeName(value.Type())}
		case *types.Func:
			result.Symbols[name] = symbol{Kind: "func", Type: signatureString(value.Type().(*types.Signature), qualifier)}
		case *types.TypeName:
			entry := symbol{Kind: "type"}
			if value.IsAlias() {
				entry.Kind = "alias"
				entry.Type = typeName(types.Unalias(value.Type()))
				result.Symbols[name] = entry
				continue
			}
			named := value.Type().(*types.Named)
			for i := 0; i < named.TypeParams().Len(); i++ {
				parameter := named.TypeParams().At(i)
				entry.TypeParameters = append(entry.TypeParameters, parameter.Obj().Name()+" "+typeName(parameter.Constraint()))
			}
			if structure, ok := named.Underlying().(*types.Struct); ok {
				entry.Type = "struct"
				for i := 0; i < structure.NumFields(); i++ {
					member := structure.Field(i)
					if member.Exported() {
						entry.Fields = append(entry.Fields, field{Name: member.Name(), Type: typeName(member.Type()), Tag: structure.Tag(i), Embedded: member.Embedded()})
					}
				}
			} else {
				entry.Type = typeName(named.Underlying())
			}
			result.Symbols[name] = entry
			for i := 0; i < named.NumMethods(); i++ {
				method := named.Method(i)
				if method.Exported() {
					signature := method.Type().(*types.Signature)
					result.Symbols[name+"."+method.Name()] = symbol{Kind: "method", Type: signatureString(signature, qualifier), Receiver: typeName(signature.Recv().Type())}
				}
			}
		}
	}
	return result, nil
}

// Parameter names do not affect compatibility; order, types and variadic shape do.
func signatureString(signature *types.Signature, qualifier types.Qualifier) string {
	name := func(typ types.Type) string { return types.TypeString(typ, qualifier) }
	var parameters, inputs, outputs []string
	for i := 0; i < signature.TypeParams().Len(); i++ {
		parameter := signature.TypeParams().At(i)
		parameters = append(parameters, parameter.Obj().Name()+" "+name(parameter.Constraint()))
	}
	for i := 0; i < signature.Params().Len(); i++ {
		typ := signature.Params().At(i).Type()
		if signature.Variadic() && i == signature.Params().Len()-1 {
			inputs = append(inputs, "..."+name(typ.(*types.Slice).Elem()))
		} else {
			inputs = append(inputs, name(typ))
		}
	}
	for i := 0; i < signature.Results().Len(); i++ {
		outputs = append(outputs, name(signature.Results().At(i).Type()))
	}
	result := "func"
	if len(parameters) > 0 {
		result += "[" + strings.Join(parameters, ", ") + "]"
	}
	result += "(" + strings.Join(inputs, ", ") + ")"
	if len(outputs) == 1 {
		result += " " + outputs[0]
	} else if len(outputs) > 1 {
		result += " (" + strings.Join(outputs, ", ") + ")"
	}
	return result
}
