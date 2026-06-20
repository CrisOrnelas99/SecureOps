import { ComponentFixture, TestBed } from '@angular/core/testing';

import { VulnerabilitiesPage } from './vulnerabilities';

describe('VulnerabilitiesPage', () => {
  let component: VulnerabilitiesPage;
  let fixture: ComponentFixture<VulnerabilitiesPage>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [VulnerabilitiesPage],
    }).compileComponents();

    fixture = TestBed.createComponent(VulnerabilitiesPage);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
