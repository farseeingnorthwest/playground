# Actors

| Evaluation    | Script      | Action          | PreLoss        | Combo              |
|---------------|-------------|-----------------|----------------|--------------------|
| `Buffer`      | `VerbActor` | `ActionBuffer`  | `LossResister` | `SelectActor`      |
| `TheoryActor` |             | `CriticalActor` | `LossStopper`  | `ProbabilityActor` |
|               |             | `ImmuneActor`   |                | `SequenceActor`    |
|               |             |                 |                | `RepeatActor`      |

# Evaluators

| Warrior          | EvaluationContext | ActionContext   |                   | Combo        |
|------------------|-------------------|-----------------|-------------------|--------------|
| `ConstEvaluator` | `SelectCounter`   | `LossEvaluator` | `CustomEvaluator` | `Adder`      |
| `AxisEvaluator`  |                   |                 |                   | `Multiplier` |
| `BuffCounter`    |                   |                 |                   |              |

# Selectors

| Side                      | Signal            | ActionSignal     | Filter               | Combo              |
|---------------------------|-------------------|------------------|----------------------|--------------------|
| `AbsoluteSideSelector`    | `CurrentSelector` | `SourceSelector` | `SortSelector`       | `PipelineSelector` |
| `SideSelector`            |                   |                  | `SuffleSelector`     |                    |
| `CounterPositionSelector` |                   |                  | `FrontSelector`      |                    |
|                           |                   |                  | `WaterLevelSelector` |                    |

# Triggers

| Signal          | Action                   | Verb                    | Combo        |
|-----------------|--------------------------|-------------------------|--------------|
| `SignalTrigger` | `CurrentIsSourceTrigger` | `VerbTrigger`           | `AnyTrigger` |
|                 | `CurrentIsTargetTrigger` | `CriticalStrikeTrigger` | `FatTrigger` |
|                 | `ReactorTrigger`         | `TagTrigger`            |              |

# Tags

| Name             | Underlying |          |
|------------------|------------|----------|
| `Label`          | `string`   |          |
| `Priority`       | `int`      |          |
| `ExclusionGroup` | `uint8`    |          |
| `StackingLimit`  | `int`      | capacity |
