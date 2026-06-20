import { Routes } from '@angular/router';
import { LoginPage } from './pages/login/login';
import { RegisterPage } from './pages/register/register';
import { DashboardPage } from './pages/dashboard/dashboard';
import { AssetsPage } from './pages/assets/assets';
import { AssetDetailsPage } from './pages/asset-details/asset-details';
import { VulnerabilitiesPage } from './pages/vulnerabilities/vulnerabilities';
import { authGuard } from './services/auth.guard';

export const routes: Routes = [
  { path: '', redirectTo: 'login', pathMatch: 'full' },
  { path: 'login', component: LoginPage },
  { path: 'register', component: RegisterPage },
  { path: 'dashboard', component: DashboardPage, canActivate: [authGuard] },
  { path: 'assets', component: AssetsPage, canActivate: [authGuard] },
  { path: 'assets/:id', component: AssetDetailsPage, canActivate: [authGuard] },
  { path: 'vulnerabilities', component: VulnerabilitiesPage, canActivate: [authGuard] },
];
