import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card, CardContent } from '../ui/card';
import { Input } from '../ui/input';
import { Badge } from '../ui/badge';
import { Search, ExternalLink } from 'lucide-react';
import { fetchRuns } from '../../lib/api';
import { Run, RunStatus } from '../../lib/types';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '../ui/table';

const statusConfig: Record<RunStatus, { label: string; className: string }> = {
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
};

export function Runs() {
  const navigate = useNavigate();
  const [search, setSearch] = useState('');
  const [runs, setRuns] = useState<Run[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchRuns()
      .then(data => {
        setRuns(Array.isArray(data) ? data : []);
        setLoading(false);
      })
      .catch(err => {
        console.error(err);
        setError('Failed to load runs');
        setRuns([]);
        setLoading(false);
      });
  }, []);

  const filteredRuns = runs.filter(run =>
    run.benchmarkName.toLowerCase().includes(search.toLowerCase()) ||
    run.agentName.toLowerCase().includes(search.toLowerCase())
  );

  if (loading) return <div className="p-8 text-center">Loading runs...</div>;
  if (error) return <div className="p-8 text-center text-red-500">{error}</div>;

  return (
    <div className="space-y-6">
      <div>
        <h1>Execuções</h1>
        <p className="text-neutral-600 dark:text-neutral-400 mt-1">
          Acompanhe o status das execuções de benchmark
        </p>
      </div>

      {/* Search */}
      <div className="relative max-w-md">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-neutral-400" />
        <Input
          placeholder="Buscar execuções..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="pl-10"
        />
      </div>

      {/* Runs Table */}
      <Card>
        <CardContent className="p-0">
          {filteredRuns.length === 0 ? (
            <div className="py-16 text-center">
              <p className="text-neutral-600 dark:text-neutral-400">
                Nenhuma execução encontrada.
              </p>
            </div>
          ) : (
            <div className="overflow-x-auto">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>ID</TableHead>
                    <TableHead>Benchmark</TableHead>
                    <TableHead>Agente</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead>Progresso</TableHead>
                    <TableHead>Início</TableHead>
                    <TableHead>Conclusão</TableHead>
                    <TableHead className="text-right">Ações</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {filteredRuns.map((run) => (
                    <TableRow
                      key={run.id}
                      className="cursor-pointer hover:bg-neutral-50 dark:hover:bg-neutral-900"
                      onClick={() => navigate(`/runs/${run.id}`)}
                    >
                      <TableCell className="font-mono">
                        {run.id}
                      </TableCell>
                      <TableCell>{run.benchmarkName}</TableCell>
                      <TableCell>{run.agentName}</TableCell>
                      <TableCell>
                        <Badge
                          variant="secondary"
                          className={statusConfig[run.status].className}
                        >
                          {statusConfig[run.status].label}
                        </Badge>
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center gap-2">
                          <div className="flex-1 h-2 bg-neutral-200 dark:bg-neutral-800 rounded-full overflow-hidden max-w-[100px]">
                            <div
                              className="h-full bg-primary transition-all"
                              style={{ width: `${run.progress}%` }}
                            />
                          </div>
                          <span className="text-neutral-600 dark:text-neutral-400">
                            {run.progress}%
                          </span>
                        </div>
                      </TableCell>
                      <TableCell>
                        {new Date(run.startedAt).toLocaleString('pt-BR', {
                          day: '2-digit',
                          month: '2-digit',
                          hour: '2-digit',
                          minute: '2-digit'
                        })}
                      </TableCell>
                      <TableCell>
                        {run.completedAt
                          ? new Date(run.completedAt).toLocaleString('pt-BR', {
                            day: '2-digit',
                            month: '2-digit',
                            hour: '2-digit',
                            minute: '2-digit'
                          })
                          : '-'}
                      </TableCell>
                      <TableCell className="text-right">
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            navigate(`/runs/${run.id}`);
                          }}
                          className="text-primary hover:underline inline-flex items-center gap-1"
                        >
                          Ver detalhes
                          <ExternalLink className="w-3 h-3" />
                        </button>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}

