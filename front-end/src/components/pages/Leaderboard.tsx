import { useState, useEffect } from 'react';
import { Card, CardContent } from '../ui/card';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '../ui/select';
import { fetchLeaderboard, fetchBenchmarks } from '../../lib/api';
import { LeaderboardEntry, Benchmark } from '../../lib/types';
import { Trophy, Medal, Award } from 'lucide-react';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '../ui/table';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '../ui/tooltip';

export function Leaderboard() {
  const [selectedBenchmark, setSelectedBenchmark] = useState<string>('');
  const [leaderboard, setLeaderboard] = useState<LeaderboardEntry[]>([]);
  const [benchmarks, setBenchmarks] = useState<Benchmark[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchBenchmarks().then(data => {
      const validData = Array.isArray(data) ? data : [];
      setBenchmarks(validData);
      if (validData.length > 0) {
        setSelectedBenchmark(validData[0].id);
      }
    }).catch(console.error);
  }, []);

  useEffect(() => {
    fetchLeaderboard().then(data => {
      setLeaderboard(Array.isArray(data) ? data : []);
      setLoading(false);
    }).catch(err => {
      console.error(err);
      setLeaderboard([]);
      setLoading(false);
    });
  }, []);

  // Filter leaderboard by selected benchmark (if backend returns all)
  // Or fetch specific leaderboard if API supports it.
  // Assuming fetchLeaderboard returns all entries for now.
  const filteredLeaderboard = leaderboard
    .filter(entry => !selectedBenchmark || entry.benchmarkId === selectedBenchmark)
    .sort((a, b) => a.rank - b.rank);

  const getRankIcon = (rank: number) => {
    switch (rank) {
      case 1:
        return <Trophy className="w-5 h-5 text-yellow-500" />;
      case 2:
        return <Medal className="w-5 h-5 text-neutral-400" />;
      case 3:
        return <Award className="w-5 h-5 text-amber-600" />;
      default:
        return <span className="w-5 h-5 flex items-center justify-center">{rank}</span>;
    }
  };

  const getPerformanceColor = (value: number, isViolation: boolean = false) => {
    if (isViolation) {
      return value === 0 ? 'text-green-600 dark:text-green-400' :
        value <= 2 ? 'text-yellow-600 dark:text-yellow-400' :
          'text-red-600 dark:text-red-400';
    }
    return value >= 90 ? 'text-green-600 dark:text-green-400' :
      value >= 80 ? 'text-yellow-600 dark:text-yellow-400' :
        'text-red-600 dark:text-red-400';
  };

  if (loading) return <div className="p-8 text-center">Loading leaderboard...</div>;

  return (
    <div className="space-y-6">
      <div>
        <h1>Leaderboard</h1>
        <p className="text-neutral-600 dark:text-neutral-400 mt-1">
          Compare o desempenho dos agentes nos benchmarks
        </p>
      </div>

      {/* Filters */}
      <div className="flex flex-col md:flex-row gap-4">
        <div className="w-full md:w-64">
          <Select value={selectedBenchmark} onValueChange={setSelectedBenchmark}>
            <SelectTrigger>
              <SelectValue placeholder="Selecione um benchmark" />
            </SelectTrigger>
            <SelectContent>
              {benchmarks.map((benchmark) => (
                <SelectItem key={benchmark.id} value={benchmark.id}>
                  {benchmark.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      </div>

      {/* Leaderboard Table */}
      <Card>
        <CardContent className="p-0">
          <div className="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="w-[60px]">Rank</TableHead>
                  <TableHead>Agente</TableHead>
                  <TableHead>
                    <TooltipProvider>
                      <Tooltip>
                        <TooltipTrigger className="cursor-help">
                          Success Rate
                        </TooltipTrigger>
                        <TooltipContent>
                          Percentual de tarefas completadas com sucesso
                        </TooltipContent>
                      </Tooltip>
                    </TooltipProvider>
                  </TableHead>
                  <TableHead>
                    <TooltipProvider>
                      <Tooltip>
                        <TooltipTrigger className="cursor-help">
                          Tool Correctness
                        </TooltipTrigger>
                        <TooltipContent>
                          Percentual de ferramentas utilizadas corretamente
                        </TooltipContent>
                      </Tooltip>
                    </TooltipProvider>
                  </TableHead>
                  <TableHead>
                    <TooltipProvider>
                      <Tooltip>
                        <TooltipTrigger className="cursor-help">
                          Violações
                        </TooltipTrigger>
                        <TooltipContent>
                          Número de restrições violadas
                        </TooltipContent>
                      </Tooltip>
                    </TooltipProvider>
                  </TableHead>
                  <TableHead>Média Turnos</TableHead>
                  <TableHead>Custo Total</TableHead>
                  <TableHead>Latência Média</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filteredLeaderboard.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={8} className="text-center py-8 text-neutral-500">
                      Nenhum dado encontrado para este benchmark.
                    </TableCell>
                  </TableRow>
                ) : (
                  filteredLeaderboard.map((entry) => (
                    <TableRow
                      key={entry.agentId}
                      className={entry.rank <= 3 ? 'bg-neutral-50 dark:bg-neutral-900/50' : ''}
                    >
                      <TableCell>
                        <div className="flex items-center justify-center">
                          {getRankIcon(entry.rank)}
                        </div>
                      </TableCell>
                      <TableCell>
                        <span className={entry.rank === 1 ? '' : ''}>
                          {entry.agentName}
                        </span>
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center gap-2">
                          <div className="flex-1 h-2 bg-neutral-200 dark:bg-neutral-800 rounded-full overflow-hidden max-w-[100px]">
                            <div
                              className="h-full bg-green-500 transition-all"
                              style={{ width: `${entry.successRate}%` }}
                            />
                          </div>
                          <span className={getPerformanceColor(entry.successRate)}>
                            {entry.successRate}%
                          </span>
                        </div>
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center gap-2">
                          <div className="flex-1 h-2 bg-neutral-200 dark:bg-neutral-800 rounded-full overflow-hidden max-w-[100px]">
                            <div
                              className="h-full bg-blue-500 transition-all"
                              style={{ width: `${entry.toolCorrectness}%` }}
                            />
                          </div>
                          <span className={getPerformanceColor(entry.toolCorrectness)}>
                            {entry.toolCorrectness}%
                          </span>
                        </div>
                      </TableCell>
                      <TableCell>
                        <span className={getPerformanceColor(entry.violations, true)}>
                          {entry.violations}
                        </span>
                      </TableCell>
                      <TableCell>{entry.avgTurns.toFixed(1)}</TableCell>
                      <TableCell>${entry.totalCost.toFixed(2)}</TableCell>
                      <TableCell>{entry.avgLatency.toFixed(1)}s</TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>

      {/* Legend */}
      <Card>
        <CardContent className="pt-6">
          <h3 className="mb-4">Legenda</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <div>
              <p className="text-neutral-600 dark:text-neutral-400">
                Success Rate
              </p>
              <p className="mt-1">
                Percentual de tarefas completadas com sucesso
              </p>
            </div>
            <div>
              <p className="text-neutral-600 dark:text-neutral-400">
                Tool Correctness
              </p>
              <p className="mt-1">
                Percentual de uso correto das ferramentas
              </p>
            </div>
            <div>
              <p className="text-neutral-600 dark:text-neutral-400">
                Violações
              </p>
              <p className="mt-1">
                Número de restrições do benchmark violadas
              </p>
            </div>
            <div>
              <p className="text-neutral-600 dark:text-neutral-400">
                Média de Turnos
              </p>
              <p className="mt-1">
                Número médio de interações por tarefa
              </p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

