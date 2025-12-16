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
        repo?: string;
        commit?: string;
        patch?: string;
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

    const handleImportSWEBench = () => {
        const input = document.createElement('input');
        input.type = 'file';
        input.accept = '.json';
        input.onchange = (e) => {
            const file = (e.target as HTMLInputElement).files?.[0];
            if (!file) return;

            const reader = new FileReader();
            reader.onload = (event) => {
                try {
                    const content = event.target?.result as string;
                    const data = JSON.parse(content);
                    const items = Array.isArray(data) ? data : [data];

                    const importedTasks = items.map((item: any) => ({
                        prompt: item.problem_statement || item.prompt || '',
                        expected_output: item.patch || item.expected_output || '', // Using patch as expected outcome/truth
                        maxTurns: 30, // Default for coding tasks
                        repo: item.repo || '',
                        commit: item.base_commit || item.commit || '',
                        patch: item.patch || '',
                        instance_id: item.instance_id || ''
                    }));

                    setTasks(importedTasks);
                    setFormData(prev => ({
                        ...prev,
                        name: prev.name || `SWE-bench Import (${items.length})`,
                        domain: 'Coding',
                        description: prev.description || 'Imported from SWE-bench dataset'
                    }));
                    toast.success(`${importedTasks.length} tasks imported!`);
                } catch (err) {
                    console.error(err);
                    toast.error('Failed to parse JSON');
                }
            };
            reader.readAsText(file);
        };
        input.click();
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        // Validate tasks
        if (tasks.some(t => !t.prompt)) {
            toast.error('Preencha os prompts das tarefas');
            return;
        }

        try {
            await createBenchmark({
                ...formData,
                tasksCount: tasks.length,
                tasks: tasks.map((t, i) => ({
                    ...t,
                    id: `task-${Date.now()}-${i}`,
                    expectedTool: t.expectedTool || '',
                    expected_output: t.expected_output || 'Check Patch', // Allow empty output if patch exists
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
            <div className="flex items-center justify-between">
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
                <Button variant="secondary" onClick={handleImportSWEBench}>
                    Import SWE-bench JSON
                </Button>
            </div>

            <div>
                <h1>Novo Benchmark</h1>
                <p className="text-neutral-600 dark:text-neutral-400 mt-1">
                    Crie um novo benchmark ou importe do SWE-bench
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
                            <CardTitle>Tarefas ({tasks.length})</CardTitle>
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
                                        <Label>Prompt / Instrução (Problem Statement)</Label>
                                        <Textarea
                                            placeholder="Digite a instrução..."
                                            value={task.prompt}
                                            onChange={(e) => handleTaskChange(index, 'prompt', e.target.value)}
                                            required
                                            className="min-h-[100px]"
                                        />
                                    </div>

                                    <details className="group">
                                        <summary className="flex items-center gap-2 cursor-pointer text-sm font-medium text-neutral-600 dark:text-neutral-400 hover:text-neutral-900 dark:hover:text-neutral-200">
                                            <span className="select-none">Detalhes & Código</span>
                                            <div className="h-px bg-neutral-200 dark:bg-neutral-800 flex-1 ml-2 group-open:bg-neutral-300 transition-all" />
                                        </summary>

                                        <div className="space-y-4 mt-4 pt-2">
                                            <div className="grid grid-cols-2 gap-4">
                                                <div className="space-y-2">
                                                    <Label>Repositório</Label>
                                                    <Input
                                                        value={task.repo || ''}
                                                        onChange={(e) => handleTaskChange(index, 'repo', e.target.value)}
                                                        placeholder="e.g. sqlfluff/sqlfluff"
                                                    />
                                                </div>
                                                <div className="space-y-2">
                                                    <Label>Base Commit</Label>
                                                    <Input
                                                        value={task.commit || ''}
                                                        onChange={(e) => handleTaskChange(index, 'commit', e.target.value)}
                                                        placeholder="SHA hash"
                                                    />
                                                </div>
                                            </div>

                                            <div className="space-y-2">
                                                <Label>Patch (Solução)</Label>
                                                <Textarea
                                                    value={task.patch || ''}
                                                    onChange={(e) => handleTaskChange(index, 'patch', e.target.value)}
                                                    placeholder="Diff content..."
                                                    className="font-mono text-xs"
                                                />
                                            </div>

                                            <div className="space-y-2">
                                                <Label>Saída Esperada (Fallback)</Label>
                                                <Textarea
                                                    placeholder="Validação textual..."
                                                    value={task.expected_output}
                                                    onChange={(e) => handleTaskChange(index, 'expected_output', e.target.value)}
                                                />
                                            </div>

                                            <div className="grid grid-cols-2 gap-4">
                                                <div className="space-y-2">
                                                    <Label>Ferramenta Esperada</Label>
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
                                                        value={task.maxTurns || 10}
                                                        onChange={(e) => handleTaskChange(index, 'maxTurns', parseInt(e.target.value))}
                                                    />
                                                </div>
                                            </div>
                                        </div>
                                    </details>
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
