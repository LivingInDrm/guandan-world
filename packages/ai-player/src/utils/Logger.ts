import * as fs from 'fs';
import * as path from 'path';

export enum LogLevel {
  ERROR = 0,
  WARN = 1,
  INFO = 2,
  DEBUG = 3,
  TRACE = 4,
}

const LEVEL_NAMES: Record<LogLevel, string> = {
  [LogLevel.ERROR]: 'ERROR',
  [LogLevel.WARN]: 'WARN',
  [LogLevel.INFO]: 'INFO',
  [LogLevel.DEBUG]: 'DEBUG',
  [LogLevel.TRACE]: 'TRACE',
};

export interface LoggerContext {
  username?: string | null;
  roomCode?: string | null;
}

export interface LoggerOptions {
  logDir?: string;
  consoleOutput?: boolean;
  fileOutput?: boolean;
  level?: LogLevel;
}

function getLogLevelFromEnv(): LogLevel {
  const envLevel = process.env.LOG_LEVEL?.toUpperCase();
  switch (envLevel) {
    case 'ERROR': return LogLevel.ERROR;
    case 'WARN': return LogLevel.WARN;
    case 'INFO': return LogLevel.INFO;
    case 'DEBUG': return LogLevel.DEBUG;
    case 'TRACE': return LogLevel.TRACE;
    default: return LogLevel.INFO;
  }
}

function formatTime(): string {
  const now = new Date();
  const hours = now.getHours().toString().padStart(2, '0');
  const minutes = now.getMinutes().toString().padStart(2, '0');
  const seconds = now.getSeconds().toString().padStart(2, '0');
  const ms = now.getMilliseconds().toString().padStart(3, '0');
  return `${hours}:${minutes}:${seconds}.${ms}`;
}

function getDateString(): string {
  const now = new Date();
  const year = now.getFullYear();
  const month = (now.getMonth() + 1).toString().padStart(2, '0');
  const day = now.getDate().toString().padStart(2, '0');
  return `${year}${month}${day}`;
}

export class Logger {
  private context: LoggerContext;
  private options: Required<LoggerOptions>;
  private currentDateString: string = '';
  private fileStream: fs.WriteStream | null = null;

  constructor(options: LoggerOptions = {}) {
    this.context = { username: null, roomCode: null };
    this.options = {
      logDir: options.logDir ?? './logs',
      consoleOutput: options.consoleOutput ?? true,
      fileOutput: options.fileOutput ?? true,
      level: options.level ?? getLogLevelFromEnv(),
    };
  }

  setContext(ctx: Partial<LoggerContext>): void {
    if (ctx.username !== undefined) {
      this.context.username = ctx.username;
    }
    if (ctx.roomCode !== undefined) {
      this.context.roomCode = ctx.roomCode;
    }
  }

  private formatContext(): string {
    const { username, roomCode } = this.context;
    if (!username) {
      return '';
    }
    if (roomCode) {
      return `[${username}@${roomCode}]`;
    }
    return `[${username}]`;
  }

  private ensureLogDir(): void {
    if (!fs.existsSync(this.options.logDir)) {
      fs.mkdirSync(this.options.logDir, { recursive: true });
    }
  }

  private getFileStream(): fs.WriteStream {
    const dateString = getDateString();
    
    if (this.currentDateString !== dateString || !this.fileStream) {
      if (this.fileStream) {
        this.fileStream.end();
      }
      this.ensureLogDir();
      const logPath = path.join(this.options.logDir, `aiplay.log.${dateString}`);
      this.fileStream = fs.createWriteStream(logPath, { flags: 'a' });
      this.fileStream.on('error', (err) => {
        console.error(`[Logger] File write error: ${err.message}`);
      });
      this.currentDateString = dateString;
    }
    
    return this.fileStream;
  }

  private log(level: LogLevel, functionName: string, message: string): void {
    if (level > this.options.level) {
      return;
    }

    const time = formatTime();
    const levelName = LEVEL_NAMES[level];
    const context = this.formatContext();
    
    const parts = [
      `[${time}]`,
      `[${levelName}]`,
      context,
      `[${functionName}]`,
      message,
    ].filter(Boolean);
    
    const line = parts.join(' ');

    if (this.options.consoleOutput) {
      if (level === LogLevel.ERROR) {
        console.error(line);
      } else if (level === LogLevel.WARN) {
        console.warn(line);
      } else {
        console.log(line);
      }
    }

    if (this.options.fileOutput) {
      const stream = this.getFileStream();
      stream.write(line + '\n');
    }
  }

  error(functionName: string, message: string): void {
    this.log(LogLevel.ERROR, functionName, message);
  }

  warn(functionName: string, message: string): void {
    this.log(LogLevel.WARN, functionName, message);
  }

  info(functionName: string, message: string): void {
    this.log(LogLevel.INFO, functionName, message);
  }

  debug(functionName: string, message: string): void {
    this.log(LogLevel.DEBUG, functionName, message);
  }

  trace(functionName: string, message: string): void {
    this.log(LogLevel.TRACE, functionName, message);
  }

  close(): void {
    if (this.fileStream) {
      this.fileStream.end();
      this.fileStream = null;
    }
  }
}

export function createLogger(options: LoggerOptions = {}): Logger {
  return new Logger(options);
}
