#!/usr/bin/env node
import { AIPlayer } from './AIPlayer.js';
import type { AIPlayerConfig } from './config.js';

function parseArgs(): AIPlayerConfig {
  const args = process.argv.slice(2);
  const config: Partial<AIPlayerConfig> = {};

  const requireValue = (arg: string, value: string | undefined): string => {
    if (!value || value.startsWith('-')) {
      console.error(`Error: ${arg} requires a value`);
      printHelp();
      process.exit(1);
    }
    return value;
  };

  for (let i = 0; i < args.length; i++) {
    const arg = args[i];
    const nextArg = args[i + 1];

    switch (arg) {
      case '--server':
      case '-s':
        config.serverUrl = requireValue(arg, nextArg);
        i++;
        break;
      case '--room-code':
      case '-r':
        config.roomCode = requireValue(arg, nextArg);
        i++;
        break;
      case '--username':
      case '-u':
        config.username = requireValue(arg, nextArg);
        i++;
        break;
      case '--password':
      case '-p':
        config.password = requireValue(arg, nextArg);
        i++;
        break;
      case '--auto-register':
      case '-a':
        config.autoRegister = true;
        break;
      case '--nickname':
      case '-n':
        config.nickname = requireValue(arg, nextArg);
        i++;
        break;
      case '--level':
      case '-l':
        config.level = parseInt(requireValue(arg, nextArg), 10);
        i++;
        break;
      case '--help':
      case '-h':
        printHelp();
        process.exit(0);
    }
  }

  if (!config.serverUrl) {
    console.error('Error: --server is required');
    printHelp();
    process.exit(1);
  }

  if (!config.roomCode) {
    console.error('Error: --room-code is required');
    printHelp();
    process.exit(1);
  }

  if (!config.autoRegister && !config.username) {
    console.error('Error: --username is required when not using --auto-register');
    printHelp();
    process.exit(1);
  }

  return config as AIPlayerConfig;
}

function printHelp(): void {
  console.log(`
Usage: npx tsx src/index.ts [options]

Options:
  -s, --server <url>      Server URL (required)
                          Example: http://localhost:8080
  -r, --room-code <code>  Room code to join (required)
                          Example: ABC123
  -u, --username <name>   Username for login
  -p, --password <pass>   Password for login/register
  -a, --auto-register     Auto-register a new user
  -n, --nickname <name>   Set player nickname
  -l, --level <num>       AI difficulty level (default: 1)
  -h, --help              Show this help message

Examples:
  # Auto-register and join room
  npx tsx src/index.ts --server http://localhost:8080 --room-code ABC123 --auto-register

  # Login with existing account
  npx tsx src/index.ts --server http://localhost:8080 --room-code ABC123 --username ai_player --password secret123
`);
}

async function main(): Promise<void> {
  const config = parseArgs();

  console.log('Starting AI Player...');
  console.log(`Server: ${config.serverUrl}`);
  console.log(`Room Code: ${config.roomCode}`);
  console.log(`Auto Register: ${config.autoRegister || false}`);

  const player = new AIPlayer(config);

  const gracefulShutdown = (signal: string) => {
    console.log(`\nReceived ${signal}, stopping...`);
    player.stop().finally(() => {
      process.exit(0);
    });
  };

  process.on('SIGINT', () => gracefulShutdown('SIGINT'));
  process.on('SIGTERM', () => gracefulShutdown('SIGTERM'));

  try {
    await player.start();
    console.log('AI Player is running. Press Ctrl+C to stop.');
  } catch (error) {
    console.error('Failed to start AI Player:', error);
    process.exit(1);
  }
}

main().catch(console.error);
