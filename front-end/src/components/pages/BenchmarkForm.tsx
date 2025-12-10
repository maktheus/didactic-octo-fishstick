import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card, CardContent, CardHeader, CardTitle } from '../ui/card';
import { Button } from '../ui/button';
import { Input } from '../ui/input';
import { Label } from '../ui/label';
import { Textarea } from '../ui/textarea';
import { ArrowLeft, Plus, Trash2 } from 'lucide-react';
import { toast } from 'sonner';
import { createBenchmark } from '../../lib/api';

export function BenchmarkForm() {
    const navigate = useNavigate();
    const [formData, setFormData] = useState({
        name: '',
        domain: '',
        description: '',
    });

    const [tasks, setTasks] = useState<{
        prompt: string;
        expected_output: string;
        expectedTool?: string;
        maxTurns?: number;
    }[]>([
        { prompt: '', expected_output: '', maxTurns: 10 }
    ]);

    const handleTaskChange = (index: number, field: string, value: string | number) => {
        const newTasks = [...tasks];
        (newTasks[index] as any)[field] = value;
        setTasks(newTasks);
    };

    const addTask = () => {
        setTasks([...tasks, { prompt: '', expected_output: '', maxTurns: 10 }]);
    };

    const removeTask = (index: number) => {
        if (tasks.length === 1) return;
        const newTasks = tasks.filter((_, i) => i !== index);
        setTasks(newTasks);
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        // Validate tasks
        if (tasks.some(t => !t.prompt || !t.expected_output)) {
            toast.error('Preencha todos os campos das tarefas');
            return;
        }

        try {
            await createBenchmark({
                ...formData,
                tasksCount: tasks.length,
                tasks: tasks.map((t, i) => ({
                    ...t,
                    id: `task-${Date.now()}-${i}`, // Generate temporary ID
                    expectedTool: t.expectedTool || '',
                    constraints: []
                }))
            });
            toast.success('Benchmark criado com sucesso!');
            setTimeout(() => {
                navigate('/benchmarks');
            }, 500);
        } catch (error) {
            console.error(error);
            toast.error('Erro ao criar benchmark');
        }
    };

    return (
        <div className="max-w-3xl space-y-6">
            <div className="flex items-center gap-4">
                <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => navigate('/benchmarks')}
                >
                    <ArrowLeft className="w-4 h-4 mr-2" />
                    Voltar
                </Button>
            </div>

            <div>
                <h1>Novo Benchmark</h1>
                <p className="text-neutral-600 dark:text-neutral-400 mt-1">
                    Crie um novo benchmark para avaliar agentes
                </p>
            </div>

            <form onSubmit={handleSubmit}>
                <div className="space-y-6">
                    <Card>
                        <CardHeader>
                            <CardTitle>Informações do Benchmark</CardTitle>
                        </CardHeader>
                        <CardContent className="space-y-6">
                            <div className="space-y-2">
                                <Label htmlFor="name">Nome do Benchmark</Label>
                                <Input
                                    id="name"
                                    placeholder="Ex: Level 4 Test"
                                    value={formData.name}
                                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                                    required
                                />
                            </div>

                            <div className="space-y-2">
                                <Label htmlFor="domain">Domínio</Label>
                                <Input
                                    id="domain"
                                    placeholder="Ex: Coding, Reasoning, Math"
                                    value={formData.domain}
                                    onChange={(e) => setFormData({ ...formData, domain: e.target.value })}
                                    required
                                />
                            </div>

                            <div className="space-y-2">
                                <Label htmlFor="description">Descrição</Label>
                                <Input
                                    id="description"
                                    placeholder="Breve descrição do benchmark"
                                    value={formData.description}
                                    onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                                    required
                                />
                            </div>
                        </CardContent>
                    </Card>

                    <Card>
                        <CardHeader className="flex flex-row items-center justify-between">
                            <CardTitle>Tarefas</CardTitle>
                            <Button type="button" variant="outline" size="sm" onClick={addTask}>
                                <Plus className="w-4 h-4 mr-2" />
                                Adicionar Tarefa
                            </Button>
                        </CardHeader>
                        <CardContent className="space-y-6">
                            {tasks.map((task, index) => (
                                <div key={index} className="p-4 border rounded-lg space-y-4 relative bg-neutral-50 dark:bg-neutral-900/50">
                                    <div className="absolute right-4 top-4">
                                        {tasks.length > 1 && (
                                            <Button
                                                type="button"
                                                variant="ghost"
                                                size="icon"
                                                className="text-red-500 hover:text-red-600 hover:bg-red-50"
                                                onClick={() => removeTask(index)}
                                            >
                                                <Trash2 className="w-4 h-4" />
                                            </Button>
                                        )}
                                    </div>

                                    <div className="space-y-2">
                                        <Label>Prompt / Instrução</Label>
                                        <Textarea
                                            placeholder="Digite a instrução para o agente..."
                                            value={task.prompt}
                                            onChange={(e) => handleTaskChange(index, 'prompt', e.target.value)}
                                            required
                                        />
                                    </div>

                                    <div className="space-y-2">
                                        <Label>Saída Esperada</Label>
                                        <Textarea
                                            placeholder="Digite a saída esperada para validação..."
                                            value={task.expected_output}
                                            onChange={(e) => handleTaskChange(index, 'expected_output', e.target.value)}
                                            required
                                        />
                                    </div>

                                    <div className="grid grid-cols-2 gap-4">
                                        <div className="space-y-2">
                                            <Label>Ferramenta Esperada (Opcional)</Label>
                                            <Input
                                                placeholder="Ex: run_command"
                                                value={task.expectedTool || ''}
                                                onChange={(e) => handleTaskChange(index, 'expectedTool', e.target.value)}
                                            />
                                        </div>
                                        <div className="space-y-2">
                                            <Label>Máximo de Turnos</Label>
                                            <Input
                                                type="number"
                                                min={1}
                                                max={50}
                                                placeholder="10"
                                                value={task.maxTurns || 10}
                                                onChange={(e) => handleTaskChange(index, 'maxTurns', parseInt(e.target.value))}
                                            />
                                        </div>
                                    </div>
                                </div>
                            ))}
                        </CardContent>
                    </Card>

                    <div className="flex gap-3 pt-4 border-t border-neutral-200 dark:border-neutral-800">
                        <Button type="submit">
                            Criar Benchmark
                        </Button>
                        <Button
                            type="button"
                            variant="outline"
                            onClick={() => navigate('/benchmarks')}
                        >
                            Cancelar
                        </Button>
                    </div>
                </div>
            </form>
        </div>
    );
}
