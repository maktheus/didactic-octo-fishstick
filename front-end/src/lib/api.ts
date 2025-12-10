const API_URL = 'http://localhost:8080';

import { Agent, Benchmark, Run, LeaderboardEntry, Trace } from './types';
import { mockAgents, mockBenchmarks, mockRuns, mockLeaderboard, mockTraces } from './mockData';

export async function fetchAgents(): Promise<Agent[]> {
    const response = await fetch(`${API_URL}/agents`);
    if (!response.ok) throw new Error('Failed to fetch agents');
    return response.json();
}

export async function fetchBenchmarks(): Promise<Benchmark[]> {
    const response = await fetch(`${API_URL}/benchmarks`);
    if (!response.ok) throw new Error('Failed to fetch benchmarks');
    return response.json();
}

export async function fetchRuns(): Promise<Run[]> {
    const response = await fetch(`${API_URL}/submissions`);
    if (!response.ok) throw new Error('Failed to fetch runs');
    const data = await response.json();
    return data.map((item: any) => {
        const scoreSummary = item.scoreSummary || {};
        return {
            ...item,
            startedAt: item.submittedAt, // Map backend field
            benchmarkName: item.benchmarkName || 'Unknown Benchmark',
            agentName: item.agentName || 'Unknown Agent',
            status: item.status || 'queued',
            // Map metrics from scoreSummary
            successRate: scoreSummary.successRate !== undefined ? scoreSummary.successRate : undefined,
            toolCorrectness: scoreSummary.toolCorrectness !== undefined ? scoreSummary.toolCorrectness : undefined,
            violations: scoreSummary.violations !== undefined ? scoreSummary.violations : undefined,
            avgTurns: scoreSummary.avgTurns !== undefined ? scoreSummary.avgTurns : undefined,
            totalCost: scoreSummary.totalCost !== undefined ? scoreSummary.totalCost : undefined,
            avgLatency: scoreSummary.avgLatency !== undefined ? scoreSummary.avgLatency : undefined
        };
    });
}

export async function fetchRun(id: string): Promise<Run> {
    const runs = await fetchRuns();
    const run = runs.find(r => r.id === id);
    if (!run) throw new Error('Run not found');
    return run;
}

export async function fetchLeaderboard(): Promise<LeaderboardEntry[]> {
    const response = await fetch(`${API_URL}/leaderboard`);
    if (!response.ok) throw new Error('Failed to fetch leaderboard');
    return response.json();
}

export async function fetchTraces(): Promise<Trace[]> {
    const response = await fetch(`${API_URL}/traces`);
    if (!response.ok) throw new Error('Failed to fetch traces');
    const events = await response.json();

    // Group events by submissionId (runId)
    const tracesMap = new Map<string, Trace>();

    events.forEach((event: any) => {
        const runId = event.submissionId;
        if (!tracesMap.has(runId)) {
            tracesMap.set(runId, {
                id: runId, // Use runId as traceId for now, assuming 1 trace per run
                runId: runId,
                taskId: event.taskId || '',
                taskName: event.taskName || 'Task Execution',
                messages: [],
                success: event.success || false,
                turns: event.turns || 0,
                cost: event.cost || 0,
                latency: event.latency || 0
            });
        }

        const trace = tracesMap.get(runId)!;

        // Map event to TraceMessage
        // Backend event types: plan, reflection, tool, etc.
        // Frontend types: user, agent, tool
        let type: 'user' | 'agent' | 'tool' = 'agent';
        if (event.type === 'tool') type = 'tool';
        if (event.type === 'user') type = 'user'; // If backend sends user messages

        trace.messages.push({
            id: event.id,
            type: type,
            content: event.message || '',
            timestamp: new Date(event.timestamp),
            toolName: event.toolName,
            parameters: event.parameters,
            result: event.result
        });

        // Update trace stats if available
        if (event.success) trace.success = true;
        if (event.turns) trace.turns = Math.max(trace.turns, event.turns);
        if (event.cost) trace.cost += event.cost; // Accumulate cost?
        if (event.latency) trace.latency += event.latency;
    });

    return Array.from(tracesMap.values());
}

export async function seedDatabase() {
    // Seed Agents
    for (const agent of mockAgents) {
        await fetch(`${API_URL}/agents`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(agent)
        });
    }

    // Seed Benchmarks
    for (const bench of mockBenchmarks) {
        await fetch(`${API_URL}/benchmarks`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(bench)
        });
    }

    // Seed Runs (Submissions)
    // Note: This might be complex if IDs need to match. 
    // For now, we'll try to post them as is, assuming backend handles ID or we use the same IDs.
    // The backend might expect a specific payload structure for submission.
    // Based on main.go: mux.HandleFunc("/submissions", ... http.MethodPost: orchestratorHTTP.Submit

    // For simplicity in this "mock to real" transition, we might skip seeding runs/traces 
    // if the backend logic is complex (e.g. requires actual execution).
    // However, we can try to inject them if the backend supports "importing" results.
    // Looking at main.go, /submissions POST triggers orchestratorHTTP.Submit which likely starts a run.
    // We probably can't easily "seed" completed runs without a dedicated import endpoint.
    // So we will only seed Agents and Benchmarks for now.

    console.log('Seeding completed for Agents and Benchmarks');
}

export async function createRun(data: { benchmarkId: string; agentId: string }) {
    const response = await fetch(`${API_URL}/submissions`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            benchmark_id: data.benchmarkId,
            agent_id: data.agentId,
            payload: "{}" // Default empty payload
        })
    });
    if (!response.ok) throw new Error('Failed to create run');
    return response.json();
}

export async function createAgent(agent: Omit<Agent, 'id' | 'createdAt' | 'status'>) {
    const response = await fetch(`${API_URL}/agents`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(agent)
    });
    if (!response.ok) throw new Error('Failed to create agent');
    return response.json();
}

export async function createBenchmark(benchmark: Omit<Benchmark, 'id' | 'createdAt'>) {
    const response = await fetch(`${API_URL}/benchmarks`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(benchmark)
    });
    if (!response.ok) throw new Error('Failed to create benchmark');
    return response.json();
}
