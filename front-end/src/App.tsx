import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { Layout } from './components/layout/Layout';
import { Dashboard } from './components/pages/Dashboard';
import { Agents } from './components/pages/Agents';
import { AgentForm } from './components/pages/AgentForm';
import { Benchmarks } from './components/pages/Benchmarks';
import { BenchmarkForm } from './components/pages/BenchmarkForm';
import { BenchmarkDetail } from './components/pages/BenchmarkDetail';
import { BenchmarkRun } from './components/pages/BenchmarkRun';
import { Runs } from './components/pages/Runs';
import { RunDetail } from './components/pages/RunDetail';
import { Leaderboard } from './components/pages/Leaderboard';
import { Traces } from './components/pages/Traces';
import { TraceViewer } from './components/pages/TraceViewer';
import { Settings } from './components/pages/Settings';
import { Toaster } from './components/ui/sonner';

export default function App() {
  return (
    <BrowserRouter>
      <Layout>
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/agents" element={<Agents />} />
          <Route path="/agents/new" element={<AgentForm />} />
          <Route path="/benchmarks" element={<Benchmarks />} />
          <Route path="/benchmarks/new" element={<BenchmarkForm />} />
          <Route path="/benchmarks/:id" element={<BenchmarkDetail />} />
          <Route path="/benchmarks/:id/run" element={<BenchmarkRun />} />
          <Route path="/runs" element={<Runs />} />
          <Route path="/runs/:id" element={<RunDetail />} />
          <Route path="/leaderboard" element={<Leaderboard />} />
          <Route path="/traces" element={<Traces />} />
          <Route path="/traces/:id" element={<TraceViewer />} />
          <Route path="/settings" element={<Settings />} />
        </Routes>
      </Layout>
      <Toaster />
    </BrowserRouter>
  );
}
