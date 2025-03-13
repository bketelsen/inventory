# https://taskfile.dev

version: '3'

tasks:
 
  site:
    desc: Run hugo dev server
    cmds:
      - task: :build
      - task: :generate
      - hugo server --buildDrafts --disableFastRender