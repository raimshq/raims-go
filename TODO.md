# TODO

- [ ] Habilitar o CI usando os checks canônicos do Developer-Toolkit depois de
  autorizar o repositório e workflow nos runners CodeBuild. Validar as versões
  suportadas de Go e o vínculo dos resultados com cada commit e release.
- [ ] Incluir SDK Go nas fontes do atendimento após confirmar a distribuição.
  Fonte responsável: guias BR/EN `integracoes/sdk-go` em raims-developer-docs,
  catálogo público consumido pelo Raims-Chat e mapa knowledge-sources.md.
  Explicar finalidade, quando sugerir, instalação e atualização, requisitos,
  acesso pela API pública, etapas, permissões da credencial, API key/OAuth,
  ambientes, disponibilidade, limites, paginação, omissão/null, idempotência,
  respostas pendentes, retries e validação de webhooks. Não é necessário mudar
  contrato de contexto nem ferramentas de consulta financeira: a feature é
  uma biblioteca cliente. Regenerar os catálogos e validar a orientação BR/EN
  antes de concluir o item; não recomendar o checkout local como pacote disponível.
- [ ] Registrar a primeira distribuição do SDK Go no changelog público em
  raims-developer-docs/content/docs/changelog/post, com data real de
  disponibilização e links para os guias. Manter pendente enquanto estiver local.
