# Test Data

Esta pasta contém dados de teste usados pelos testes da biblioteca.

## Estrutura

- `events/` - Eventos JSON reais de produção organizados por categoria:
  - `ack-events/` - Eventos de confirmação de ordens e uploads (3 arquivos)
  - `dms-events/` - Eventos de Driver Monitoring System (11 arquivos)
  - `driver-behavior-events/` - Eventos de comportamento do motorista (7 arquivos)
  - `hardware-events/` - Eventos de hardware e sistema (14 arquivos)
  - `telemetry-events/` - Eventos de telemetria e localização (6 arquivos)
  - `vision-basic-events/` - Eventos básicos de visão (5 arquivos)

Total: 46 arquivos JSON com exemplos reais de eventos.

## Uso nos Testes

Os testes em `internal/parser/parser_test.go` usam esses arquivos para validar o parsing de eventos reais.
