import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '../ui/card';
import { Label } from '../ui/label';
import { Input } from '../ui/input';
import { Switch } from '../ui/switch';
import { Button } from '../ui/button';
import { toast } from 'sonner';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '../ui/dialog';

export function Settings() {
  const [resetPassword, setResetPassword] = useState('');
  const [isDialogOpen, setIsDialogOpen] = useState(false);

  const handleSave = () => {
    toast.success('Configurações salvas com sucesso!');
  };

  const handleReset = () => {
    if (resetPassword === 'w1ntersun') {
      import('../../lib/api').then(({ resetPlatform }) => {
        toast.promise(resetPlatform(), {
          loading: 'Resetando plataforma...',
          success: 'Plataforma resetada com sucesso!',
          error: 'Erro ao resetar plataforma'
        });
      });
      setIsDialogOpen(false);
      setResetPassword('');
    } else {
      toast.error('Senha incorreta!');
    }
  };

  return (
    <div className="max-w-4xl space-y-6">
      <div>
        <h1>Configurações</h1>
        <p className="text-neutral-600 dark:text-neutral-400 mt-1">
          Gerencie as configurações da plataforma
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Perfil</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="name">Nome</Label>
            <Input id="name" defaultValue="Administrador" />
          </div>
          <div className="space-y-2">
            <Label htmlFor="email">Email</Label>
            <Input id="email" type="email" defaultValue="admin@agentbench.com" />
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Preferências</CardTitle>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="flex items-center justify-between">
            <div className="space-y-0.5">
              <Label>Notificações por Email</Label>
              <p className="text-neutral-600 dark:text-neutral-400">
                Receba notificações quando execuções forem concluídas
              </p>
            </div>
            <Switch defaultChecked />
          </div>

          <div className="flex items-center justify-between">
            <div className="space-y-0.5">
              <Label>Tema Escuro</Label>
              <p className="text-neutral-600 dark:text-neutral-400">
                Ativar modo escuro da interface
              </p>
            </div>
            <Switch />
          </div>

          <div className="flex items-center justify-between">
            <div className="space-y-0.5">
              <Label>Auto-executar Benchmarks</Label>
              <p className="text-neutral-600 dark:text-neutral-400">
                Executar automaticamente novos benchmarks quando criados
              </p>
            </div>
            <Switch />
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>API</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="api-key">API Key</Label>
            <div className="flex gap-2">
              <Input
                id="api-key"
                type="password"
                defaultValue="••••••••••••••••••••"
                className="flex-1"
              />
              <Button variant="outline">Regenerar</Button>
            </div>
            <p className="text-neutral-500">
              Use esta chave para acessar a API do AgentBench
            </p>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Dados de Teste</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label>Popular Banco de Dados</Label>
            <p className="text-neutral-600 dark:text-neutral-400">
              Adiciona dados de exemplo (Agentes e Benchmarks) para testes.
            </p>
            <Button
              variant="outline"
              onClick={() => {
                import('../../lib/api').then(({ seedDatabase }) => {
                  toast.promise(seedDatabase(), {
                    loading: 'Populando banco de dados...',
                    success: 'Banco de dados populado com sucesso!',
                    error: 'Erro ao popular banco de dados'
                  });
                });
              }}
            >
              Popular Dados
            </Button>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Zona de Perigo</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label>Limpar Dados</Label>
            <p className="text-neutral-600 dark:text-neutral-400">
              Remove todas as execuções e traces, mantendo agentes e benchmarks
            </p>
            <Button variant="outline" className="text-red-600 hover:text-red-700">
              Limpar Execuções
            </Button>
          </div>

          <div className="space-y-2 pt-4 border-t border-neutral-200 dark:border-neutral-800">
            <Label>Resetar Plataforma</Label>
            <p className="text-neutral-600 dark:text-neutral-400">
              Remove todos os dados da plataforma. Esta ação não pode ser desfeita.
            </p>

            <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
              <DialogTrigger asChild>
                <Button variant="outline" className="text-red-600 hover:text-red-700">
                  Resetar Tudo
                </Button>
              </DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>Resetar Plataforma</DialogTitle>
                  <DialogDescription>
                    Tem certeza que deseja apagar TODOS os dados? Esta ação é irreversível.
                    Digite a senha de administrador para confirmar.
                  </DialogDescription>
                </DialogHeader>
                <div className="space-y-4 py-4">
                  <div className="space-y-2">
                    <Label htmlFor="admin-password">Senha de Administrador</Label>
                    <Input
                      id="admin-password"
                      type="password"
                      placeholder="Digite a senha..."
                      value={resetPassword}
                      onChange={(e) => setResetPassword(e.target.value)}
                    />
                  </div>
                  <Button
                    variant="destructive"
                    className="w-full"
                    onClick={handleReset}
                  >
                    Confirmar Reset
                  </Button>
                </div>
              </DialogContent>
            </Dialog>
          </div>
        </CardContent>
      </Card>

      <div className="flex gap-3">
        <Button onClick={handleSave}>Salvar Configurações</Button>
        <Button variant="outline">Cancelar</Button>
      </div>
    </div>
  );
}
