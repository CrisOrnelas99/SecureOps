import { Routes } from '@angular/router';
import { Login } from "./pages/login/login";
import { Register } from "./pages/register/register";
import { Dashboard } from "./pages/dashboard/dashboard";
import { Assets } from "./pages/assets/assets";
import { AssetDetails } from './pages/asset-details/asset-details';
import { Vulnerabilities} from './pages/vulnerabilities/vulnerabilities';

export const routes: Routes = [
  { path: '', redirectTo: 'login', pathMatch: 'full' },
  { path: 'login', component: Login },
  { path: 'register', component: Register },
  { path: 'dashboard', component: Dashboard },
  { path: 'assets', component: Assets },
  { path: 'assets/:id', component: AssetDetails },
  { path: 'vulnerabilities', component: Vulnerabilities },
];
