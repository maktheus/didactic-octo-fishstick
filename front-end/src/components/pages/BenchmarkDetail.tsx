import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Card, CardContent, CardHeader, CardTitle } from '../ui/card';
import { Button } from '../ui/button';
import { Badge } from '../ui/badge';
import { ArrowLeft, PlayCircle } from 'lucide-react';
import { fetchBenchmarks } from '../../lib/api';
import { Benchmark } from '../../lib/types';

export function BenchmarkDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [benchmark, setBenchmark] = useState<Benchmark | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchBenchmarks()
      .then(data => {
        const found = data.find(b => b.id === id);
        setBenchmark(found || null);
        setLoading(false);
      })
      .catch(err => {
        console.error(err);
        setLoading(false);
      });
  }, [id]);

  if (loading) {
    return <div className="p-8 text-center">Carregando...</div>;
  }

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
            <p className="mt-2">{benchmark.tasks ? benchmark.tasks.length : 0}</p>
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
              {benchmark.createdAt ? new Date(benchmark.createdAt).toLocaleDateString('pt-BR') : '-'}
            </p>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Tarefas do Benchmark</CardTitle>
        </CardHeader>
        <CardContent>
          {!benchmark.tasks || benchmark.tasks.length === 0 ? (
            <p className="text-neutral-600 dark:text-neutral-400 text-center py-8">
              Nenhuma tarefa cadastrada ainda.
            </p>
          ) : (
            <div className="space-y-4">
              {benchmark.tasks.map((task, index) => (
                <div
                  key={index}
                  className="p-4 border border-neutral-200 dark:border-neutral-800 rounded-lg"
                >
                  <div className="flex items-start gap-4">
                    <div className="flex items-center justify-center w-8 h-8 rounded-full bg-primary/10 text-primary flex-shrink-0">
                      {index + 1}
                    </div>
                    <div className="flex-1 space-y-2">
                      <p>{task.prompt}</p>

                      {task.expected_output && (
                        <div className="flex items-center gap-2">
                          <span className="text-neutral-600 dark:text-neutral-400">
                            Saída Esperada:
                          </span>
                          <span className="font-mono text-sm bg-neutral-100 dark:bg-neutral-800 px-2 py-1 rounded">
                            {task.expected_output}
                          </span>
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
