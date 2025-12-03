import { useState, useEffect } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '../ui/card';
import { Bot, Target, PlayCircle, TrendingUp } from 'lucide-react';
import { fetchAgents, fetchBenchmarks, fetchRuns } from '../../lib/api';
import { Agent, Benchmark, Run } from '../../lib/types';
import { BarChart, Bar, LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';

export function Dashboard() {
  const [agents, setAgents] = useState<Agent[]>([]);
  const [benchmarks, setBenchmarks] = useState<Benchmark[]>([]);
  const [runs, setRuns] = useState<Run[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    Promise.all([
      fetchAgents(),
      fetchBenchmarks(),
      fetchRuns()
    ]).then(([agentsData, benchmarksData, runsData]) => {
      setAgents(Array.isArray(agentsData) ? agentsData : []);
      setBenchmarks(Array.isArray(benchmarksData) ? benchmarksData : []);
      setRuns(Array.isArray(runsData) ? runsData : []);
      setLoading(false);
    }).catch(err => {
      console.error('Failed to load dashboard data', err);
      setAgents([]);
      setBenchmarks([]);
      setRuns([]);
      setLoading(false);
    });
  }, []);

  if (loading) return <div className="p-8 text-center">Loading dashboard...</div>;

  const activeAgents = agents.filter(a => a.status === 'active').length;
  const totalBenchmarks = benchmarks.length;
  const completedRuns = runs.filter(r => r.status === 'completed').length;
  const runningRuns = runs.filter(r => r.status === 'running').length;

  // Success rate by agent (mock calculation for now as we need more data)
  // In a real app, we'd aggregate run results by agent
  const successByAgent = agents.slice(0, 5).map(agent => {
    const agentRuns = runs.filter(r => r.agentId === agent.id && r.status === 'completed');
    const successRate = agentRuns.length > 0
      ? agentRuns.reduce((acc, r) => acc + (r.successRate || 0), 0) / agentRuns.length
      : 0;
    return { name: agent.name, rate: Math.round(successRate * 100) };
  });

  // Runs over time (mock aggregation)
  const runsOverTime = [
    { date: '05/11', runs: 3 },
    { date: '06/11', runs: 5 },
    { date: '07/11', runs: 8 },
    { date: '08/11', runs: 6 },
    { date: '09/11', runs: 10 },
    { date: '10/11', runs: 7 },
  ];

  const stats = [
    { name: 'Agentes Ativos', value: activeAgents, icon: Bot, color: 'text-blue-600' },
    { name: 'Benchmarks', value: totalBenchmarks, icon: Target, color: 'text-purple-600' },
    { name: 'Execuções Concluídas', value: completedRuns, icon: PlayCircle, color: 'text-green-600' },
    { name: 'Em Execução', value: runningRuns, icon: TrendingUp, color: 'text-orange-600' },
  ];

  return (
    <div className="space-y-6">
      <div>
        <h1>Dashboard</h1>
        <p className="text-neutral-600 dark:text-neutral-400 mt-1">
          Visão geral da plataforma de benchmark
        </p>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {stats.map((stat) => (
          <Card key={stat.name}>
            <CardContent className="pt-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-neutral-600 dark:text-neutral-400">
                    {stat.name}
                  </p>
                  <p className="mt-2">{stat.value}</p>
                </div>
                <stat.icon className={`w-8 h-8 ${stat.color}`} />
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Charts */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <Card>
          <CardHeader>
            <CardTitle>Taxa de Sucesso por Agente</CardTitle>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={280}>
              <BarChart data={successByAgent}>
                <CartesianGrid strokeDasharray="3 3" stroke="#e5e5e5" />
                <XAxis dataKey="name" stroke="#737373" />
                <YAxis stroke="#737373" />
                <Tooltip
                  contentStyle={{
                    backgroundColor: '#fff',
                    border: '1px solid #e5e5e5',
                    borderRadius: '8px'
                  }}
                />
                <Bar dataKey="rate" fill="hsl(var(--primary))" radius={[4, 4, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Execuções ao Longo do Tempo</CardTitle>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={280}>
              <LineChart data={runsOverTime}>
                <CartesianGrid strokeDasharray="3 3" stroke="#e5e5e5" />
                <XAxis dataKey="date" stroke="#737373" />
                <YAxis stroke="#737373" />
                <Tooltip
                  contentStyle={{
                    backgroundColor: '#fff',
                    border: '1px solid #e5e5e5',
                    borderRadius: '8px'
                  }}
                />
                <Line
                  type="monotone"
                  dataKey="runs"
                  stroke="hsl(var(--primary))"
                  strokeWidth={2}
                  dot={{ fill: 'hsl(var(--primary))' }}
                />
              </LineChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>
      </div>

      {/* Recent Runs */}
      <Card>
        <CardHeader>
          <CardTitle>Execuções Recentes</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {runs.slice(0, 5).map((run) => (
              <div
                key={run.id}
                className="flex items-center justify-between p-4 border border-neutral-200 dark:border-neutral-800 rounded-lg"
              >
                <div className="flex-1">
                  <p>{run.benchmarkName}</p>
                  <p className="text-neutral-600 dark:text-neutral-400 mt-1">
                    {run.agentName}
                  </p>
                </div>
                <div className="flex items-center gap-4">
                  <div className="text-right">
                    <p className="text-neutral-600 dark:text-neutral-400">
                      Progresso
                    </p>
                    <p className="mt-1">{run.progress}%</p>
                  </div>
                  <div className={`px-3 py-1 rounded-full ${run.status === 'completed' ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400' :
                    run.status === 'running' ? 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400' :
                      run.status === 'failed' ? 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400' :
                        'bg-neutral-100 text-neutral-700 dark:bg-neutral-800 dark:text-neutral-400'
                    }`}>
                    {run.status === 'completed' ? 'Concluído' :
                      run.status === 'running' ? 'Em execução' :
                        run.status === 'failed' ? 'Falhou' : 'Pendente'}
                  </div>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

