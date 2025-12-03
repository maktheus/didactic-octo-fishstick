import { useParams, useNavigate } from 'react-router-dom';
import { Card, CardContent } from '../ui/card';
import { Button } from '../ui/button';
import { Badge } from '../ui/badge';
import { ArrowLeft, User, Bot, Wrench } from 'lucide-react';
import { mockTraces } from '../../lib/mockData';

export function TraceViewer() {
  const { id } = useParams();
  const navigate = useNavigate();
  const trace = mockTraces.find(t => t.id === id);

  if (!trace) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <Button variant="ghost" size="sm" onClick={() => navigate('/traces')}>
            <ArrowLeft className="w-4 h-4 mr-2" />
            Voltar
          </Button>
        </div>
        <Card>
          <CardContent className="py-16 text-center">
            <p className="text-neutral-600 dark:text-neutral-400">
              Trace não encontrado.
            </p>
          </CardContent>
        </Card>
      </div>
    );
  }

  const getMessageIcon = (type: string) => {
    switch (type) {
      case 'user':
        return <User className="w-5 h-5" />;
      case 'agent':
        return <Bot className="w-5 h-5" />;
      case 'tool':
        return <Wrench className="w-5 h-5" />;
      default:
        return null;
    }
  };

  const getMessageColor = (type: string) => {
    switch (type) {
      case 'user':
        return 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400';
      case 'agent':
        return 'bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-400';
      case 'tool':
        return 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400';
      default:
        return 'bg-neutral-100 text-neutral-700 dark:bg-neutral-800 dark:text-neutral-400';
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Button variant="ghost" size="sm" onClick={() => navigate('/traces')}>
          <ArrowLeft className="w-4 h-4 mr-2" />
          Voltar
        </Button>
      </div>

      <div>
        <div className="flex items-center gap-3">
          <h1>{trace.taskName}</h1>
          <Badge variant={trace.success ? 'default' : 'secondary'}>
            {trace.success ? 'Sucesso' : 'Falha'}
          </Badge>
        </div>
        <p className="text-neutral-600 dark:text-neutral-400 mt-2">
          Trace ID: {trace.id} • Run ID: {trace.runId}
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="pt-6">
            <p className="text-neutral-600 dark:text-neutral-400">
              Turnos
            </p>
            <p className="mt-2">{trace.turns}</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-6">
            <p className="text-neutral-600 dark:text-neutral-400">
              Custo
            </p>
            <p className="mt-2">
              ${trace.cost.toFixed(2)}
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-6">
            <p className="text-neutral-600 dark:text-neutral-400">
              Latência
            </p>
            <p className="mt-2">
              {trace.latency.toFixed(1)}s
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-6">
            <p className="text-neutral-600 dark:text-neutral-400">
              Mensagens
            </p>
            <p className="mt-2">{trace.messages.length}</p>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardContent className="pt-6">
          <h2 className="mb-6">Timeline de Interações</h2>
          <div className="space-y-4">
            {trace.messages.map((message, index) => (
              <div key={message.id} className="flex gap-4">
                <div className="flex flex-col items-center">
                  <div className={`w-10 h-10 rounded-full flex items-center justify-center ${getMessageColor(message.type)}`}>
                    {getMessageIcon(message.type)}
                  </div>
                  {index < trace.messages.length - 1 && (
                    <div className="w-0.5 flex-1 bg-neutral-200 dark:bg-neutral-800 my-2" />
                  )}
                </div>

                <div className="flex-1 pb-8">
                  <div className="flex items-center gap-2 mb-2">
                    <span className="capitalize">
                      {message.type === 'user' ? 'Usuário' : 
                       message.type === 'agent' ? 'Agente' : 
                       'Ferramenta'}
                    </span>
                    <span className="text-neutral-500">•</span>
                    <span className="text-neutral-600 dark:text-neutral-400">
                      {message.timestamp.toLocaleTimeString('pt-BR')}
                    </span>
                  </div>

                  <Card>
                    <CardContent className="pt-4">
                      <p className="whitespace-pre-wrap">{message.content}</p>

                      {message.toolName && (
                        <div className="mt-4 p-3 bg-neutral-50 dark:bg-neutral-900 rounded-lg">
                          <p className="text-neutral-600 dark:text-neutral-400 mb-2">
                            Ferramenta: <span className="font-mono">{message.toolName}</span>
                          </p>
                          
                          {message.parameters && (
                            <div className="mt-2">
                              <p className="text-neutral-600 dark:text-neutral-400 mb-1">
                                Parâmetros:
                              </p>
                              <pre className="text-sm bg-neutral-100 dark:bg-neutral-950 p-3 rounded overflow-x-auto">
                                {JSON.stringify(message.parameters, null, 2)}
                              </pre>
                            </div>
                          )}

                          {message.result && (
                            <div className="mt-2">
                              <p className="text-neutral-600 dark:text-neutral-400 mb-1">
                                Resultado:
                              </p>
                              <pre className="text-sm bg-neutral-100 dark:bg-neutral-950 p-3 rounded overflow-x-auto">
                                {JSON.stringify(message.result, null, 2)}
                              </pre>
                            </div>
                          )}
                        </div>
                      )}
                    </CardContent>
                  </Card>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
