import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Card, CardContent, CardHeader, CardTitle } from '../ui/card';
import { Button } from '../ui/button';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '../ui/select';
import { Label } from '../ui/label';
import { fetchAgents, fetchBenchmarks, createRun } from '../../lib/api';
import { Agent, Benchmark } from '../../lib/types';
import { toast } from 'sonner';

export function BenchmarkRun() {
    const { id } = useParams();
    const navigate = useNavigate();
    const [agents, setAgents] = useState<Agent[]>([]);
    const [benchmark, setBenchmark] = useState<Benchmark | null>(null);
    const [selectedAgent, setSelectedAgent] = useState<string>('');
    const [loading, setLoading] = useState(true);
    const [submitting, setSubmitting] = useState(false);

    useEffect(() => {
        Promise.all([
            fetchAgents(),
            fetchBenchmarks()
        ]).then(([agentsData, benchmarksData]) => {
            setAgents(Array.isArray(agentsData) ? agentsData : []);
            const foundBenchmark = benchmarksData.find(b => b.id === id);
            setBenchmark(foundBenchmark || null);
            setLoading(false);
        }).catch(err => {
            console.error(err);
            toast.error('Erro ao carregar dados');
            setLoading(false);
        });
    }, [id]);

    const handleRun = async () => {
        if (!selectedAgent || !benchmark) return;

        setSubmitting(true);
        try {
            await createRun({
                benchmarkId: benchmark.id,
                agentId: selectedAgent
            });
            toast.success('Execução iniciada com sucesso!');
            navigate('/runs');
        } catch (error) {
            console.error(error);
            toast.error('Erro ao iniciar execução');
        } finally {
            setSubmitting(false);
        }
    };

    if (loading) return <div className="p-8 text-center">Carregando...</div>;
    if (!benchmark) return <div className="p-8 text-center">Benchmark não encontrado</div>;

    return (
        <div className="max-w-2xl mx-auto space-y-6">
            <div>
                <h1>Executar Benchmark</h1>
                <p className="text-neutral-600 dark:text-neutral-400 mt-1">
                    Configure e inicie uma nova execução
                </p>
            </div>

            <Card>
                <CardHeader>
                    <CardTitle>{benchmark.name}</CardTitle>
                </CardHeader>
                <CardContent className="space-y-6">
                    <div className="space-y-2">
                        <Label>Benchmark</Label>
                        <div className="p-3 bg-neutral-100 dark:bg-neutral-800 rounded-md">
                            {benchmark.description}
                        </div>
                    </div>

                    <div className="space-y-2">
                        <Label>Selecione o Agente</Label>
                        <Select value={selectedAgent} onValueChange={setSelectedAgent}>
                            <SelectTrigger>
                                <SelectValue placeholder="Selecione um agente..." />
                            </SelectTrigger>
                            <SelectContent>
                                {agents.map((agent) => (
                                    <SelectItem key={agent.id} value={agent.id}>
                                        {agent.name} ({agent.provider})
                                    </SelectItem>
                                ))}
                            </SelectContent>
                        </Select>
                    </div>

                    <div className="flex justify-end gap-3 pt-4">
                        <Button variant="outline" onClick={() => navigate('/benchmarks')}>
                            Cancelar
                        </Button>
                        <Button onClick={handleRun} disabled={!selectedAgent || submitting}>
                            {submitting ? 'Iniciando...' : 'Iniciar Execução'}
                        </Button>
                    </div>
                </CardContent>
            </Card>
        </div>
    );
}
