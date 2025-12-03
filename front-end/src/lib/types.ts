export interface Agent {
  id: string;
  name: string;
  provider: string;
  endpoint: string;
  model?: string;
  systemPrompt?: string;
  authType: 'none' | 'bearer' | 'apikey';
  authToken?: string;
  status: 'active' | 'inactive';
  createdAt: string;
  headers?: Record<string, string>;
}

export interface Task {
  id: string;
  prompt: string;
  expectedTool?: string;
  constraints?: string[];
  maxTurns?: number;
}

export interface Benchmark {
  id: string;
  name: string;
  description: string;
  domain: string;
  tasks: Task[];
  createdAt: Date;
  tasksCount: number;
}

export type RunStatus = 'pending' | 'running' | 'completed' | 'failed';

export interface Run {
  id: string;
  benchmarkId: string;
  benchmarkName: string;
  agentId: string;
  agentName: string;
  status: RunStatus;
  startedAt: Date;
  completedAt?: Date;
  progress: number;
  successRate?: number;
  toolCorrectness?: number;
  violations?: number;
  avgTurns?: number;
  totalCost?: number;
  avgLatency?: number;
}

export interface LeaderboardEntry {
  rank: number;
  agentId: string;
  agentName: string;
  successRate: number;
  toolCorrectness: number;
  violations: number;
  avgTurns: number;
  totalCost: number;
  avgLatency: number;
  benchmarkId: string;
}

export interface TraceMessage {
  id: string;
  type: 'user' | 'agent' | 'tool';
  content: string;
  timestamp: Date;
  toolName?: string;
  parameters?: Record<string, any>;
  result?: any;
}

export interface Trace {
  id: string;
  runId: string;
  taskId: string;
  taskName: string;
  messages: TraceMessage[];
  success: boolean;
  turns: number;
  cost: number;
  latency: number;
}
