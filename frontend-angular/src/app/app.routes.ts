import { Routes } from '@angular/router';
import { Login } from './pages/login/login';
import { Register } from './pages/register/register';
import { Dashboard } from './pages/dashboard/dashboard';
import { Assets } from './pages/assets/assets';
import { AssetDetails } from './pages/asset-details/asset-details';
import { Vulnerabilities } from './pages/vulnerabilities/vulnerabilities';
import { authGuard } from './services/auth.guard';

export const routes: Routes = [
  { path: '', redirectTo: 'login', pathMatch: 'full' },
  { path: 'login', component: Login },
  { path: 'register', component: Register },
  { path: 'dashboard', component: Dashboard, canActivate: [authGuard] },
  { path: 'assets', component: Assets, canActivate: [authGuard] },
  { path: 'assets/:id', component: AssetDetails, canActivate: [authGuard] },
  { path: 'vulnerabilities', component: Vulnerabilities, canActivate: [authGuard] },
];
