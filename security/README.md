# Verificação estática do SDK Go

O preflight usa Semgrep 1.181.0 e as regras oficiais na revisão
`9bbf021d0acffa3095fcffd17b2bc7787eb4796d`, junto das regras locais em
`security/semgrep.yml`.

A análise cobre toda a raiz do SDK, incluindo os arquivos Go gerados e os
scripts de engenharia em Go e Python. O resultado precisa comprovar a leitura
de cada fonte `.go` e `.py`; somente os diretórios de artefatos previstos pelo
validador canônico são excluídos.

As regras locais conferem verificação TLS, configuração isolada dos clientes
HTTP, execução de shells e credenciais incorporadas. As regras oficiais de Go
e Python complementam essa análise. Erros do scanner, escopo incompleto ou
achados não resolvidos bloqueiam a validação.

O helper canônico requer Python 3.12. No exemplo abaixo, `python` precisa
apontar para essa versão. No Windows, o launcher `Test-RaimsRelease.ps1`
também aceita `-PythonPath` para selecionar o executável.

O comando pode ser executado no workspace:

```powershell
python Developer-Toolkit/scripts/release-preflight/sdk_security.py --source Raims-Go --workspace . --repository Raims-Go --output .tmp/raims-go-sdk/security --prepare-tools
```

Use uma saída nova ou vazia. Os relatórios ficam somente em `.tmp`. Essa
verificação não cria infraestrutura, não habilita workflows e não publica o SDK.
