import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Card, CardContent } from '../ui/card';
import { Button } from '../ui/button';
import { Badge } from '../ui/badge';
import { ArrowLeft, Activity } from 'lucide-react';
import { fetchRun, fetchTraces, fetchBenchmarks } from '../../lib/api';
import { Run, Trace, Benchmark } from '../../lib/types';
import { Progress } from '../ui/progress';

export function RunDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [run, setRun] = useState<Run | null>(null);
  const [traces, setTraces] = useState<Trace[]>([]);
  const [loading, setLoading] = useState(true);

  const [benchmark, setBenchmark] = useState<Benchmark | null>(null);

  useEffect(() => {
    if (id) {
      Promise.all([
        fetchRun(id).catch(() => null),
        fetchTraces().then(ts => ts.filter(t => t.runId === id)).catch(() => [])
      ]).then(async ([runData, tracesData]) => {
        setRun(runData);
        setTraces(tracesData);

        if (runData && runData.benchmarkId) {
          try {
            const allBenchs = await fetchBenchmarks();
            const bench = allBenchs.find(b => b.id === runData.benchmarkId);
            setBenchmark(bench || null);
          } catch (e) {
            console.error(e);
          }
        }
        setLoading(false);
      });
    }
  }, [id]);

  if (loading) {
    return <div className="p-8 text-center">Carregando detalhes...</div>;
  }

  if (!run) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <Button variant="ghost" size="sm" onClick={() => navigate('/runs')}>
            <ArrowLeft className="w-4 h-4 mr-2" />
            Voltar
          </Button>
        </div>
        <Card>
          <CardContent className="py-16 text-center">
            <p className="text-neutral-600 dark:text-neutral-400">
              Execução não encontrada.
            </p>
          </CardContent>
        </Card>
      </div>
    );
  }

  const statusConfig: Record<string, { label: string; className: string }> = {
    completed: {
      label: 'Concluído',
      className: 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400'
    },
    running: {
      label: 'Em execução',
      className: 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400'
    },
    failed: {
      label: 'Falhou',
      className: 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400'
    },
    pending: {
      label: 'Pendente',
      className: 'bg-neutral-100 text-neutral-700 dark:bg-neutral-800 dark:text-neutral-400'
    },
    queued: {
      label: 'Na Fila',
      className: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400'
    },
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Button variant="ghost" size="sm" onClick={() => navigate('/runs')}>
          <ArrowLeft className="w-4 h-4 mr-2" />
          Voltar
        </Button>
      </div>

      <div className="flex flex-col md:flex-row md:items-start md:justify-between gap-4">
        <div>
          <div className="flex items-center gap-3">
            <h1>Execução {run.id}</h1>
            <Badge
              variant="secondary"
              className={statusConfig[run.status]?.className || statusConfig['pending'].className}
            >
              {statusConfig[run.status]?.label || run.status}
            </Badge>
          </div>
          <p className="text-neutral-600 dark:text-neutral-400 mt-2 flex items-center gap-2 flex-wrap">
            {run.benchmarkName} • {run.agentName}
            {benchmark && benchmark.tasks?.[0]?.repo && (
              <span className="bg-neutral-100 dark:bg-neutral-800 px-2 py-0.5 rounded text-xs font-mono text-neutral-500">
                📦 {benchmark.tasks[0].repo}
              </span>
            )}
            {benchmark && benchmark.tasks?.[0]?.commit && (
              <span className="bg-neutral-100 dark:bg-neutral-800 px-2 py-0.5 rounded text-xs font-mono text-neutral-500">
                🔗 {benchmark.tasks[0].commit.substring(0, 8)}
              </span>
            )}
          </p>
        </div>
        {run.status === 'completed' && (
          <Button onClick={() => navigate(`/traces?runId=${run.id}`)}>
            <Activity className="w-4 h-4 mr-2" />
            Ver Traces
          </Button>
        )}
      </div>

      <Card>
        <CardContent className="pt-6">
          <div className="space-y-4">
            <div>
              <div className="flex items-center justify-between mb-2">
                <span className="text-neutral-600 dark:text-neutral-400">
                  Progresso
                </span>
                <span>{run.progress}%</span>
              </div>
              <Progress value={run.progress} />
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 pt-4 border-t border-neutral-200 dark:border-neutral-800">
              <div>
                <p className="text-neutral-600 dark:text-neutral-400">
                  Início
                </p>
                <p className="mt-1">
                  {run.startedAt.toLocaleString('pt-BR')}
                </p>
              </div>
              {run.completedAt && (
                <div>
                  <p className="text-neutral-600 dark:text-neutral-400">
                    Conclusão
                  </p>
                  <p className="mt-1">
                    {run.completedAt.toLocaleString('pt-BR')}
                  </p>
                </div>
              )}
            </div>
          </div>
        </CardContent>

      </Card>

      {/* Plan Display - Commented out as Trace type structure changed
      {
        traces.find(t => t.type === 'plan') && (
          <Card>
            <CardContent className="pt-6">
              <h2 className="mb-4 flex items-center gap-2">
                <Activity className="w-4 h-4" />
                Plano de Execução
              </h2>
              <div className="bg-neutral-50 dark:bg-neutral-900 p-4 rounded-lg font-mono text-sm whitespace-pre-wrap">
                {traces.find(t => t.type === 'plan')?.message}
              </div>
            </CardContent>
          </Card>
        )
      }
      */}

      {
        run.status === 'completed' && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <Card>
              <CardContent className="pt-6">
                <p className="text-neutral-600 dark:text-neutral-400">
                  Taxa de Sucesso
                </p>
                <p className="mt-2">{run.successRate}%</p>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="pt-6">
                <p className="text-neutral-600 dark:text-neutral-400">
                  Tool Correctness
                </p>
                <p className="mt-2">{run.toolCorrectness}%</p>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="pt-6">
                <p className="text-neutral-600 dark:text-neutral-400">
                  Violações
                </p>
                <p className="mt-2">{run.violations}</p>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="pt-6">
                <p className="text-neutral-600 dark:text-neutral-400">
                  Média de Turnos
                </p>
                <p className="mt-2">{run.avgTurns}</p>
              </CardContent>
            </Card>
          </div>
        )
      }

      {
        run.totalCost !== undefined && run.avgLatency !== undefined && (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Card>
              <CardContent className="pt-6">
                <p className="text-neutral-600 dark:text-neutral-400">
                  Custo Total
                </p>
                <p className="mt-2">
                  ${run.totalCost.toFixed(2)}
                </p>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="pt-6">
                <p className="text-neutral-600 dark:text-neutral-400">
                  Latência Média
                </p>
                <p className="mt-2">
                  {run.avgLatency.toFixed(1)}s
                </p>
              </CardContent>
            </Card>
          </div>
        )
      }

      {
        traces.length > 0 && (
          <Card>
            <CardContent className="pt-6">
              <h2 className="mb-4">Traces</h2>
              <div className="space-y-2">
                {traces.map((trace) => (
                  <button
                    key={trace.id}
                    onClick={() => navigate(`/traces/${trace.id}`)}
                    className="w-full p-4 border border-neutral-200 dark:border-neutral-800 rounded-lg hover:bg-neutral-50 dark:hover:bg-neutral-900 transition-colors text-left"
                  >
                    <div className="flex items-center justify-between">
                      <div>
                        <p>{trace.taskName}</p>
                        <p className="text-neutral-600 dark:text-neutral-400 mt-1">
                          {trace.turns} turnos • {trace.latency}s
                        </p>
                      </div>
                      <Badge variant={trace.success ? 'default' : 'secondary'}>
                        {trace.success ? 'Sucesso' : 'Falha'}
                      </Badge>
                    </div>
                  </button>
                ))}
              </div>
            </CardContent>
          </Card>
        )
      }
    </div >
  );
}
