// Browser entry point for the Angular application.
import { bootstrapApplication } from '@angular/platform-browser';
import { appConfig } from './app/config/app.config';
import { App } from './app/app';

bootstrapApplication(App, appConfig).catch((err) => console.error(err));
