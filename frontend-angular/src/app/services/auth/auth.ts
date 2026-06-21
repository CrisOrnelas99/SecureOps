// Client-side authentication state used by the route guard.
import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  // Backend authentication is not wired yet, so this remains a placeholder gate.
  isAuthenticated(): boolean {
    return false;
  }
}
