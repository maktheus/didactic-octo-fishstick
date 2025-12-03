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
    // The backend might return runs in a different format or endpoint
    // Assuming /results or /submissions based on main.go
    // main.go has /results -> runnerHTTP.Results and /submissions -> orchestratorHTTP.List
    // Let's use /submissions for now as it seems to map to "Runs"
    const response = await fetch(`${API_URL}/submissions`);
    if (!response.ok) throw new Error('Failed to fetch runs');
    return response.json();
}

export async function fetchLeaderboard(): Promise<LeaderboardEntry[]> {
    const response = await fetch(`${API_URL}/leaderboard`);
    if (!response.ok) throw new Error('Failed to fetch leaderboard');
    return response.json();
}

export async function fetchTraces(): Promise<Trace[]> {
    const response = await fetch(`${API_URL}/traces`);
    if (!response.ok) throw new Error('Failed to fetch traces');
    return response.json();
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
