import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AssetDetailsPage } from './asset-details';

describe('AssetDetailsPage', () => {
  let component: AssetDetailsPage;
  let fixture: ComponentFixture<AssetDetailsPage>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AssetDetailsPage],
    }).compileComponents();

    fixture = TestBed.createComponent(AssetDetailsPage);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
