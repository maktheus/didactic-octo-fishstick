import { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Card, CardContent } from '../ui/card';
import { Button } from '../ui/button';
import { Input } from '../ui/input';
import { Badge } from '../ui/badge';
import { Plus, Search, PlayCircle } from 'lucide-react';
import { fetchBenchmarks } from '../../lib/api';
import { Benchmark } from '../../lib/types';

export function Benchmarks() {
  const navigate = useNavigate();
  const [search, setSearch] = useState('');
  const [benchmarks, setBenchmarks] = useState<Benchmark[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchBenchmarks()
      .then(data => {
        setBenchmarks(Array.isArray(data) ? data : []);
        setLoading(false);
      })
      .catch(err => {
        console.error(err);
        setError('Failed to load benchmarks');
        setBenchmarks([]);
        setLoading(false);
      });
  }, []);

  const filteredBenchmarks = benchmarks.filter(benchmark =>
    benchmark.name.toLowerCase().includes(search.toLowerCase()) ||
    benchmark.domain.toLowerCase().includes(search.toLowerCase())
  );

  if (loading) return <div className="p-8 text-center">Loading benchmarks...</div>;
  if (error) return <div className="p-8 text-center text-red-500">{error}</div>;

  return (
    <div className="space-y-6">
      <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
        <div>
          <h1>Benchmarks</h1>
          <p className="text-neutral-600 dark:text-neutral-400 mt-1">
            Gerencie os benchmarks e cenários de teste
          </p>
        </div>
        <Button>
          <Plus className="w-4 h-4 mr-2" />
          Novo Benchmark
        </Button>
      </div>

      {/* Search */}
      <div className="relative max-w-md">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-neutral-400" />
        <Input
          placeholder="Buscar benchmarks..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="pl-10"
        />
      </div>

      {/* Benchmarks List */}
      {filteredBenchmarks.length === 0 ? (
        <Card>
          <CardContent className="py-16 text-center">
            <p className="text-neutral-600 dark:text-neutral-400">
              Nenhum benchmark encontrado.
            </p>
            <Button variant="outline" className="mt-4">
              <Plus className="w-4 h-4 mr-2" />
              Criar Primeiro Benchmark
            </Button>
          </CardContent>
        </Card>
      ) : (
        <div className="space-y-4">
          {filteredBenchmarks.map((benchmark) => (
            <Card
              key={benchmark.id}
              className="hover:shadow-md transition-shadow cursor-pointer"
              onClick={() => navigate(`/benchmarks/${benchmark.id}`)}
            >
              <CardContent className="pt-6">
                <div className="flex flex-col md:flex-row md:items-start md:justify-between gap-4">
                  <div className="flex-1">
                    <div className="flex items-center gap-3 flex-wrap">
                      <h2>{benchmark.name}</h2>
                      <Badge variant="outline">{benchmark.domain}</Badge>
                    </div>
                    <p className="text-neutral-600 dark:text-neutral-400 mt-2">
                      {benchmark.description}
                    </p>
                    <div className="flex items-center gap-6 mt-4">
                      <div>
                        <span className="text-neutral-600 dark:text-neutral-400">
                          Tarefas:
                        </span>
                        <span className="ml-2">{benchmark.tasksCount}</span>
                      </div>
                      <div>
                        <span className="text-neutral-600 dark:text-neutral-400">
                          Criado em:
                        </span>
                        <span className="ml-2">
                          {new Date(benchmark.createdAt).toLocaleDateString('pt-BR')}
                        </span>
                      </div>
                    </div>
                  </div>
                  <Button
                    onClick={(e) => {
                      e.stopPropagation();
                      navigate(`/benchmarks/${benchmark.id}/run`);
                    }}
                  >
                    <PlayCircle className="w-4 h-4 mr-2" />
                    Executar
                  </Button>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}

