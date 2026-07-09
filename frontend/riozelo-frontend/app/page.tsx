'use client';

import { useEffect, useState } from 'react';

export default function TesteAmbiente() {
  const [logSSE, setLogSSE] = useState<string[]>([]);

  // RF04 - Escuta o servidor Go em tempo real
  useEffect(() => {
    const eventSource = new EventSource('http://localhost:8080/dashboard/stream');

    eventSource.onmessage = (event) => {
      setLogSSE((prev) => [`Nova ocorrência recebida: ${event.data}`, ...prev]);
    };

    return () => eventSource.close();
  }, []);

  // Função para simular o clique de envio do formulário
  const enviarDadosTeste = async () => {
    const dadosFake = {
      macroCategoria: "Vias e Asfalto",
      subcategoria: "Buraco na Rua",
      detalheOutros: "",
      bairro: "Madureira",
      rua: "Estrada do Portela"
    };

    await fetch('http://localhost:8080/ocorrencias', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(dadosFake)
    });
  };

  return (
    <div style={{ padding: '20px', fontFamily: 'sans-serif', background: '#111', color: '#fff', minHeight: '100vh' }}>
      <h1>RioZelo — Teste de Conexão (MVP 1) 🌴</h1>
      <hr style={{ borderColor: '#333' }} />

      <div style={{ marginTop: '20px' }}>
        <h3>1. Lado do Cidadão (Simular Envio)</h3>
        <button onClick={enviarDadosTeste} style={{ padding: '10px 20px', background: '#0070f3', color: '#fff', border: 'none', borderRadius: '5px', cursor: 'pointer' }}>
          Disparar Ocorrência de Teste
        </button>
      </div>

      <div style={{ marginTop: '40px' }}>
        <h3>2. Lado do Servidor (Log do Canal SSE em Tempo Real)</h3>
        <div style={{ background: '#222', padding: '15px', borderRadius: '5px', minHeight: '150px', fontSize: '12px', fontFamily: 'monospace' }}>
          {logSSE.length === 0 && <p style={{ color: '#666' }}>Aguardando cliques ou eventos...</p>}
          {logSSE.map((log, index) => (
            <p key={index} style={{ color: '#4af626', margin: '5px 0' }}>{log}</p>
          ))}
        </div>
      </div>
    </div>
  );
}