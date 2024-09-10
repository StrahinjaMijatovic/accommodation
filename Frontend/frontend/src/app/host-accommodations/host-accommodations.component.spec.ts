import { ComponentFixture, TestBed } from '@angular/core/testing';

import { HostAccommodationsComponent } from './host-accommodations.component';

describe('HostAccommodationsComponent', () => {
  let component: HostAccommodationsComponent;
  let fixture: ComponentFixture<HostAccommodationsComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [HostAccommodationsComponent]
    });
    fixture = TestBed.createComponent(HostAccommodationsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
