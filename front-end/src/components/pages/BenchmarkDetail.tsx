import { useParams, useNavigate } from 'react-router-dom';
import { Card, CardContent, CardHeader, CardTitle } from '../ui/card';
import { Button } from '../ui/button';
import { Badge } from '../ui/badge';
import { ArrowLeft, PlayCircle } from 'lucide-react';
import { mockBenchmarks } from '../../lib/mockData';

export function BenchmarkDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const benchmark = mockBenchmarks.find(b => b.id === id);

  if (!benchmark) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <Button variant="ghost" size="sm" onClick={() => navigate('/benchmarks')}>
            <ArrowLeft className="w-4 h-4 mr-2" />
            Voltar
          </Button>
        </div>
        <Card>
          <CardContent className="py-16 text-center">
            <p className="text-neutral-600 dark:text-neutral-400">
              Benchmark não encontrado.
            </p>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Button variant="ghost" size="sm" onClick={() => navigate('/benchmarks')}>
          <ArrowLeft className="w-4 h-4 mr-2" />
          Voltar
        </Button>
      </div>

      <div className="flex flex-col md:flex-row md:items-start md:justify-between gap-4">
        <div>
          <div className="flex items-center gap-3">
            <h1>{benchmark.name}</h1>
            <Badge variant="outline">{benchmark.domain}</Badge>
          </div>
          <p className="text-neutral-600 dark:text-neutral-400 mt-2">
            {benchmark.description}
          </p>
        </div>
        <Button onClick={() => navigate(`/benchmarks/${benchmark.id}/run`)}>
          <PlayCircle className="w-4 h-4 mr-2" />
          Executar Benchmark
        </Button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card>
          <CardContent className="pt-6">
            <p className="text-neutral-600 dark:text-neutral-400">
              Total de Tarefas
            </p>
            <p className="mt-2">{benchmark.tasksCount}</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-6">
            <p className="text-neutral-600 dark:text-neutral-400">
              Domínio
            </p>
            <p className="mt-2">{benchmark.domain}</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-6">
            <p className="text-neutral-600 dark:text-neutral-400">
              Criado em
            </p>
            <p className="mt-2">
              {benchmark.createdAt.toLocaleDateString('pt-BR')}
            </p>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Tarefas do Benchmark</CardTitle>
        </CardHeader>
        <CardContent>
          {benchmark.tasks.length === 0 ? (
            <p className="text-neutral-600 dark:text-neutral-400 text-center py-8">
              Nenhuma tarefa cadastrada ainda.
            </p>
          ) : (
            <div className="space-y-4">
              {benchmark.tasks.map((task, index) => (
                <div
                  key={task.id}
                  className="p-4 border border-neutral-200 dark:border-neutral-800 rounded-lg"
                >
                  <div className="flex items-start gap-4">
                    <div className="flex items-center justify-center w-8 h-8 rounded-full bg-primary/10 text-primary flex-shrink-0">
                      {index + 1}
                    </div>
                    <div className="flex-1 space-y-2">
                      <p>{task.prompt}</p>
                      
                      {task.expectedTool && (
                        <div className="flex items-center gap-2">
                          <span className="text-neutral-600 dark:text-neutral-400">
                            Ferramenta esperada:
                          </span>
                          <Badge variant="secondary">{task.expectedTool}</Badge>
                        </div>
                      )}
                      
                      {task.constraints && task.constraints.length > 0 && (
                        <div>
                          <span className="text-neutral-600 dark:text-neutral-400">
                            Restrições:
                          </span>
                          <ul className="mt-2 space-y-1 ml-4">
                            {task.constraints.map((constraint, i) => (
                              <li key={i} className="text-neutral-600 dark:text-neutral-400 list-disc">
                                {constraint}
                              </li>
                            ))}
                          </ul>
                        </div>
                      )}
                      
                      {task.maxTurns && (
                        <div className="flex items-center gap-2">
                          <span className="text-neutral-600 dark:text-neutral-400">
                            Turnos máximos:
                          </span>
                          <span>{task.maxTurns}</span>
                        </div>
                      )}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
