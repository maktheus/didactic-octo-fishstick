import { useState, useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { Card, CardContent } from '../ui/card';
import { Badge } from '../ui/badge';
import { fetchTraces } from '../../lib/api';
import { Trace } from '../../lib/types';

export function Traces() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const runId = searchParams.get('runId');
  const [traces, setTraces] = useState<Trace[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchTraces().then(data => {
      setTraces(Array.isArray(data) ? data : []);
      setLoading(false);
    }).catch(err => {
      console.error(err);
      setTraces([]);
      setLoading(false);
    });
  }, []);

  const filteredTraces = runId
    ? traces.filter(t => t.runId === runId)
    : traces;

  if (loading) return <div className="p-8 text-center">Loading traces...</div>;

  return (
    <div className="space-y-6">
      <div>
        <h1>Traces</h1>
        <p className="text-neutral-600 dark:text-neutral-400 mt-1">
          Visualize os traces detalhados das execuções
        </p>
      </div>

      {filteredTraces.length === 0 ? (
        <Card>
          <CardContent className="py-16 text-center">
            <p className="text-neutral-600 dark:text-neutral-400">
              Nenhum trace encontrado.
            </p>
          </CardContent>
        </Card>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {filteredTraces.map((trace) => (
            <Card
              key={trace.id}
              className="hover:shadow-md transition-shadow cursor-pointer"
              onClick={() => navigate(`/traces/${trace.id}`)}
            >
              <CardContent className="pt-6">
                <div className="flex items-start justify-between mb-3">
                  <h3>{trace.taskName}</h3>
                  <Badge variant={trace.success ? 'default' : 'secondary'}>
                    {trace.success ? 'Sucesso' : 'Falha'}
                  </Badge>
                </div>

                <div className="space-y-2">
                  <div className="flex justify-between">
                    <span className="text-neutral-600 dark:text-neutral-400">
                      Run ID
                    </span>
                    <span className="font-mono">{trace.runId}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-neutral-600 dark:text-neutral-400">
                      Turnos
                    </span>
                    <span>{trace.turns}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-neutral-600 dark:text-neutral-400">
                      Mensagens
                    </span>
                    <span>{trace.messages.length}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-neutral-600 dark:text-neutral-400">
                      Latência
                    </span>
                    <span>{trace.latency.toFixed(1)}s</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-neutral-600 dark:text-neutral-400">
                      Custo
                    </span>
                    <span>${trace.cost.toFixed(2)}</span>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}

