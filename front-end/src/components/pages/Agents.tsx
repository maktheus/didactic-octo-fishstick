import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { Card, CardContent } from '../ui/card';
import { Button } from '../ui/button';
import { Input } from '../ui/input';
import { Badge } from '../ui/badge';
import { Plus, Search, MoreVertical, Bot, Globe } from 'lucide-react';
import { fetchAgents } from '../../lib/api';
import { Agent } from '../../lib/types';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '../ui/dropdown-menu';

export function Agents() {
  const [search, setSearch] = useState('');
  const [agents, setAgents] = useState<Agent[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchAgents()
      .then(data => {
        setAgents(Array.isArray(data) ? data : []);
        setLoading(false);
      })
      .catch(err => {
        console.error(err);
        setError('Failed to load agents');
        setAgents([]);
        setLoading(false);
      });
  }, []);

  const filteredAgents = agents.filter(agent =>
    agent.name.toLowerCase().includes(search.toLowerCase()) ||
    agent.provider.toLowerCase().includes(search.toLowerCase())
  );

  if (loading) return <div className="p-8 text-center">Loading agents...</div>;
  if (error) return <div className="p-8 text-center text-red-500">{error}</div>;

  return (
    <div className="space-y-6">
      <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
        <div>
          <h1>Agentes</h1>
          <p className="text-neutral-600 dark:text-neutral-400 mt-1">
            Gerencie os agentes de IA cadastrados na plataforma
          </p>
        </div>
        <Link to="/agents/new">
          <Button>
            <Plus className="w-4 h-4 mr-2" />
            Novo Agente
          </Button>
        </Link>
      </div>

      {/* Search */}
      <div className="relative max-w-md">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-neutral-400" />
        <Input
          placeholder="Buscar agentes..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="pl-10"
        />
      </div>

      {/* Agents List */}
      {filteredAgents.length === 0 ? (
        <Card>
          <CardContent className="py-16 text-center">
            <p className="text-neutral-600 dark:text-neutral-400">
              Nenhum agente encontrado.
            </p>
            <Link to="/agents/new">
              <Button variant="outline" className="mt-4">
                <Plus className="w-4 h-4 mr-2" />
                Criar Primeiro Agente
              </Button>
            </Link>
          </CardContent>
        </Card>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {filteredAgents.map((agent) => (
            <Card key={agent.id} className="hover:shadow-md transition-shadow">
              <CardContent className="pt-6">
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <div className="flex items-center gap-2">
                      <h3>{agent.name}</h3>
                      <Badge variant={agent.status === 'active' ? 'default' : 'secondary'}>
                        {agent.status === 'active' ? 'Ativo' : 'Inativo'}
                      </Badge>
                    </div>
                    <p className="text-neutral-600 dark:text-neutral-400 mt-1">
                      {agent.provider}
                    </p>
                  </div>
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="ghost" size="sm">
                        <MoreVertical className="w-4 h-4" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      <DropdownMenuItem>Editar</DropdownMenuItem>
                      <DropdownMenuItem>Ver Execuções</DropdownMenuItem>
                      <DropdownMenuItem className="text-red-600">
                        Excluir
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </div>

                <div className="flex justify-between">
                  <span className="text-neutral-600 dark:text-neutral-400">
                    Autenticação
                  </span>
                  <span className="capitalize">{agent.authType}</span>
                </div>

                <div className="pt-4 border-t border-neutral-200 dark:border-neutral-800 space-y-2">
                  <div className="flex items-center gap-2 text-sm text-neutral-500">
                    <Bot className="w-4 h-4" />
                    <span>{agent.provider}</span>
                    {agent.model && (
                      <>
                        <span>•</span>
                        <span>{agent.model}</span>
                      </>
                    )}
                  </div>
                  <div className="flex items-center gap-2 text-sm text-neutral-500">
                    <Globe className="w-4 h-4" />
                    <span className="truncate max-w-[200px]">{agent.endpoint}</span>
                  </div>
                  {agent.systemPrompt && (
                    <div className="text-xs text-neutral-400 mt-2 line-clamp-2 italic">
                      "{agent.systemPrompt}"
                    </div>
                  )}
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}

