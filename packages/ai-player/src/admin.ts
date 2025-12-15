#!/usr/bin/env node
import express from 'express';
import cors from 'cors';
import path from 'path';
import { fileURLToPath } from 'url';
import { PlayerManager } from './admin/PlayerManager.js';
import { createRouter } from './admin/routes.js';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

function parseArgs(): { port: number } {
  const args = process.argv.slice(2);
  let port = 3001;

  for (let i = 0; i < args.length; i++) {
    if (args[i] === '--port' && args[i + 1]) {
      port = parseInt(args[i + 1], 10);
      i++;
    }
  }

  return { port };
}

async function main(): Promise<void> {
  const { port } = parseArgs();
  const manager = new PlayerManager();

  const app = express();
  app.use(cors());
  app.use(express.json());

  app.use('/api', createRouter(manager));

  const adminUiPath = path.join(__dirname, '../dist/admin-ui');
  app.use(express.static(adminUiPath));
  app.get('*', (_req, res) => {
    res.sendFile(path.join(adminUiPath, 'index.html'));
  });

  const server = app.listen(port, () => {
    console.log(`Admin server running at http://localhost:${port}`);
  });

  const gracefulShutdown = async (signal: string) => {
    console.log(`\nReceived ${signal}, stopping all players...`);
    const shutdownTimeout = setTimeout(() => {
      console.error('Shutdown timeout, forcing exit');
      process.exit(1);
    }, 10000);
    try {
      await manager.stopAll();
      server.close(() => {
        clearTimeout(shutdownTimeout);
        console.log('Server closed');
        process.exit(0);
      });
    } catch (error) {
      console.error('Error during shutdown:', error);
      clearTimeout(shutdownTimeout);
      process.exit(1);
    }
  };

  process.on('SIGINT', () => gracefulShutdown('SIGINT'));
  process.on('SIGTERM', () => gracefulShutdown('SIGTERM'));
}

main().catch(console.error);
