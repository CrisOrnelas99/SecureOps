// Server rendering rules for the Angular routes.
import { RenderMode, ServerRoute } from '@angular/ssr';

export const serverRoutes: ServerRoute[] = [
  {
    path: 'dashboard',
    renderMode: RenderMode.Client,
  },
  {
    path: 'assets',
    renderMode: RenderMode.Client,
  },
  {
    path: 'assets/:id',
    renderMode: RenderMode.Client,
  },
  {
    path: 'vulnerabilities',
    renderMode: RenderMode.Client,
  },
  {
    path: '**',
    renderMode: RenderMode.Prerender,
  },
];
