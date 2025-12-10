# Evolução da Arquitetura para Benchmark de Agentes (Agent-Centric)

Este documento detalha as mudanças necessárias para transformar o atual *Model Benchmark* (onde o sistema controla o loop do agente) em um verdadeiro *Agent Benchmark* (como SWE-bench), onde o agente é uma entidade independente avaliada por sua capacidade de resolver tarefas.

## 1. Conceito Fundamental

**Atual (Model-Centric)**:
- O Runner roda o loop: *Plan -> Execute -> Reflect*.
- O Runner chama a API do OpenAI (o "cérebro").
- O Sistema define COMO o agente trabalha.

**Alvo (Agent-Centric)**:
- O Agente é uma "Caixa Preta" (Black Box).
- O Runner apenas fornece o **Ambiente** e a **Tarefa**.
- O Agente decide seu próprio loop (ReAct, COT, Plan-and-Solve, etc.).

---

## 2. Mudanças Arquiteturais Necessárias

### A. Protocolo de Interface de Agente (Agent-Protocol)
Precisamos definir um contrato padrão que todo Agente submetido deve implementar. Isso desacopla o Runner da implementação do Agente.

**Exemplo de API (REST/HTTP) que o Agente deve expor:**

```json
POST /agent/step
Request: {
  "task_id": "123",
  "observation": "Output do comando anterior ou descrição inicial...",
  "tools_available": [...]
}
Response: {
  "action": "run_command",
  "action_input": "ls -la",
  "thought": "Preciso listar os arquivos para entender o diretório."
}
```

### B. Containerização de Agentes
Em vez de configurar apenas `Model + Endpoint`, o usuário deve submeter uma **Imagem Docker** do agente.

**Fluxo de Execução**:
1. **Runner** inicia o container do **Ambiente** (onde o código vulnerável está).
2. **Runner** inicia o container do **Agente** (a imagem submetida pelo usuário).
3. **Network Bridge**: Runner conecta os dois containers ou atua como proxy.

### C. Refatoração do Runner (`runner_service.go`)

O Runner deixa de ser o "Executor" e passa a ser o "Orquestrador".

**Novo Loop do Runner:**
1. Inicializa Sandbox (Environment).
2. Inicializa Agente (Container).
3. **Loop Principal**:
    - Captura estado do Sandbox (stdout, stderr, arquivos).
    - Envia estado para o Agente (`POST /step`).
    - Recebe Ação do Agente.
    - Executa Ação no Sandbox.
    - Se Ação == `Submit` ou Terminar, encerra e avalia.

---

## 3. Plano de Implementação Sugerido

### Fase 1: Padronização da API (The "Harness")
Crie um servidor wrapper em Python/Go que implemente o loop atual, mas exponha a API proposta acima. Isso permite testar a nova arquitetura sem mudar os agentes atuais imediatamente.

### Fase 2: Suporte a Containers de Agentes
Altere o `AgentService` para aceitar `ImageURL` em vez de apenas `ModelName`.
- Use a API do Docker para subir o container do agente ao lado do sandbox.

### Fase 3: Benchmark "Headless"
Permita que o agente tenha acesso total ao shell, sem restrição de tools específicas (como o SWE-bench faz). O agente deve ser capaz de instalar dependências e rodar scripts livremente dentro do container de ambiente.

## 4. Exemplo de Definição de Agente (YAML)

Em vez de apenas selecionar "GPT-4", uma submissão de agente seria:

```yaml
agent:
  name: "SuperCoder-v1"
  image: "registry.example.com/supercoder:latest"
  resources:
    cpu: "2"
    memory: "4Gi"
  env:
    OPENAI_API_KEY: "${SECRET_KEY}"
```

## 5. Comparativo Visual

| Feature | Atual (V1) | Proposto (V2) |
| :--- | :--- | :--- |
| **Foco** | Avaliar LLMs (Models) | Avaliar Sistemas Autônomos (Agents) |
| **Controle** | Runner dita o fluxo | Agente dita o fluxo |
| **Extensibilidade** | Baixa (Hardcoded Loop) | Alta (Qualquer container Docker) |
| **Compatibilidade** | Apenas OpenAI-compatível | Qualquer linguagem/framework |

---

## 6. Próximos Passos Imediatos

Se você deseja seguir para esse nível:

1. [ ] Definir a especificação "Agent-Protocol" (v1).
2. [ ] Criar um "Reference Agent" que implementa esse protocolo (pode ser o atual encapsulado).
3. [ ] Alterar o Runner para falar com esse protocolo via HTTP em vez de chamar OpenAI diretamente.

## 7. Evolução do Frontend

Para suportar essa nova arquitetura, o Frontend precisará de ajustes para permitir a submissão de agentes e a visualização de execuções agnósticas.

### A. Submissão de Agentes (Novo Form)
Atualmente, o sistema assume que o "Agente" é apenas uma config. Precisaremos de uma tela de **"New Agent Submission"**:
- **Image Repository**: `registry.docker.com/my-agent:v1`
- **Resources**: CPU/RAM limits.
- **Environment**: Variáveis de ambiente secretas (API Keys) injetadas no container.

### B. Visualização de Steps (TraceViewer)
O `TraceViewer.tsx` atual renderiza mensagens de User/Agent/Tool. Ele precisará ser adaptado para o novo **Agent-Protocol**:
- Em vez de assumir que toda ferramenta é uma função OpenAI, ele deve renderizar a **Ação Bruta** (ex: `Action: run_command("ls -la")`).
- **Logs do Container**: Adicionar uma aba ou painel para ver o `stdout/stderr` cru do container do agente, vital para debugging de crashes.

### C. Run Detail
- **Remover dependência de "Plan"**: O componente `RunDetail.tsx` atual busca hardcoded um trace do tipo `plan`. Isso deve ser opcional, pois nem todo agente faz planejamento explícito.
- **Métricas Personalizadas**: O frontend deve ser capaz de renderizar métricas dinâmicas retornadas pelo benchmark, não apenas "Score" e "Accuracy".
