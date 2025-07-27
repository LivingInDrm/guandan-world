import '@testing-library/jest-dom';

// Mock environment variables
Object.defineProperty(import.meta, 'env', {
  value: {
    VITE_API_BASE_URL: 'http://localhost:8080'
  }
});

// Mock fetch for API tests
global.fetch = vi.fn();