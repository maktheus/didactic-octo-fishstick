import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card, CardContent, CardHeader, CardTitle } from '../ui/card';
import { Button } from '../ui/button';
import { Input } from '../ui/input';
import { Label } from '../ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '../ui/select';
import { ArrowLeft } from 'lucide-react';
import { toast } from 'sonner';
import { createAgent } from '../../lib/api';

export function AgentForm() {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    name: '',
    provider: '',
    endpoint: '',
    image: '', // New
    model: '',
    systemPrompt: '',
    authType: 'none' as 'none' | 'bearer' | 'apikey',
    authToken: '',
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      await createAgent({
        ...formData,
        // headers: {} // Add if needed
      });
      toast.success('Agente criado com sucesso!');
      setTimeout(() => {
        navigate('/agents');
      }, 500);
    } catch (error) {
      console.error(error);
      toast.error('Erro ao criar agente');
    }
  };

  return (
    <div className="max-w-3xl space-y-6">
      <div className="flex items-center gap-4">
        <Button
          variant="ghost"
          size="sm"
          onClick={() => navigate('/agents')}
        >
          <ArrowLeft className="w-4 h-4 mr-2" />
          Voltar
        </Button>
      </div>

      <div>
        <h1>Novo Agente</h1>
        <p className="text-neutral-600 dark:text-neutral-400 mt-1">
          Cadastre um novo agente de IA na plataforma
        </p>
      </div>

      <form onSubmit={handleSubmit}>
        <Card>
          <CardHeader>
            <CardTitle>Informações do Agente</CardTitle>
          </CardHeader>
          <CardContent className="space-y-6">
            <div className="space-y-2">
              <Label htmlFor="name">Nome do Agente</Label>
              <Input
                id="name"
                placeholder="Ex: GPT-4 Agent"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                required
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="provider">Provedor</Label>
              <Select
                value={formData.provider}
                onValueChange={(value: string) => {
                  let newEndpoint = formData.endpoint;
                  if (value === 'OpenAI') newEndpoint = 'https://api.openai.com/v1/chat/completions';
                  setFormData({ ...formData, provider: value, endpoint: newEndpoint });
                }}
              >
                <SelectTrigger id="provider">
                  <SelectValue placeholder="Selecione..." />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="OpenAI">OpenAI (GPT-4)</SelectItem>
                  <SelectItem value="Anthropic">Anthropic (Claude)</SelectItem>
                  <SelectItem value="Ollama">Ollama (Local)</SelectItem>
                  <SelectItem value="Custom">Custom / Agent Protocol</SelectItem>
                </SelectContent>
              </Select>
            </div>

            {(formData.provider === 'Custom' || formData.provider === 'Ollama') && (
              <div className="space-y-2">
                <Label htmlFor="endpoint">Endpoint URL</Label>
                <Input
                  id="endpoint"
                  type="url"
                  placeholder={formData.provider === 'Ollama' ? "http://host.docker.internal:11434/v1/chat/completions" : "https://api.example.com/v1/chat"}
                  value={formData.endpoint}
                  onChange={(e) => setFormData({ ...formData, endpoint: e.target.value })}
                  required
                />
                <p className="text-neutral-500 text-sm">
                  {formData.provider === 'Ollama' ? 'Use host.docker.internal para acessar o Ollama rodando no host.' : 'URL completa do endpoint da API.'}
                </p>
              </div>
            )}

            <div className="space-y-2">
              <Label htmlFor="image">Docker Image (Opcional)</Label>
              <Input
                id="image"
                placeholder="Ex: ref-agent:latest"
                value={formData.image}
                onChange={(e) => setFormData({ ...formData, image: e.target.value })}
              />
              <p className="text-neutral-500 text-sm">
                Se preenchido, o sistema iniciará um container. Deixe em branco para usar APIs (OpenAI/Ollama).
              </p>
            </div>

            <div className="space-y-2">
              <Label htmlFor="model">Modelo</Label>
              <Input
                id="model"
                placeholder="Ex: gpt-4, llama3"
                value={formData.model}
                onChange={(e) => setFormData({ ...formData, model: e.target.value })}
                required
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="systemPrompt">System Prompt (Instruções)</Label>
              <textarea
                id="systemPrompt"
                className="flex min-h-[80px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                placeholder="Você é um assistente útil..."
                value={formData.systemPrompt}
                onChange={(e) => setFormData({ ...formData, systemPrompt: e.target.value })}
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="authType">Autenticação</Label>
              <Select
                value={formData.authType}
                onValueChange={(value: 'none' | 'bearer' | 'apikey') =>
                  setFormData({ ...formData, authType: value })
                }
              >
                <SelectTrigger id="authType">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="none">Nenhuma (Mock/Local)</SelectItem>
                  <SelectItem value="bearer">Bearer Token (OpenAI)</SelectItem>
                  <SelectItem value="apikey">API Key (Genérica)</SelectItem>
                </SelectContent>
              </Select>
            </div>

            {formData.authType !== 'none' && (
              <div className="space-y-2">
                <Label htmlFor="authToken">
                  {formData.authType === 'bearer' ? 'Bearer Token (sk-...)' : 'Chave de API'}
                </Label>
                <Input
                  id="authToken"
                  type="password"
                  placeholder="••••••••••••••••"
                  value={formData.authToken}
                  onChange={(e) => setFormData({ ...formData, authToken: e.target.value })}
                  required
                />
              </div>
            )}

            <div className="flex gap-3 pt-4 border-t border-neutral-200 dark:border-neutral-800">
              <Button type="submit">
                Criar Agente
              </Button>
              <Button
                type="button"
                variant="outline"
                onClick={() => navigate('/agents')}
              >
                Cancelar
              </Button>
            </div>
          </CardContent>
        </Card>
      </form>
    </div>
  );
}
