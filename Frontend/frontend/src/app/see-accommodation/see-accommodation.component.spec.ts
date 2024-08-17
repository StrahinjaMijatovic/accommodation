import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SeeAccommodationComponent } from './see-accommodation.component';

describe('SeeAccommodationComponent', () => {
  let component: SeeAccommodationComponent;
  let fixture: ComponentFixture<SeeAccommodationComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [SeeAccommodationComponent]
    });
    fixture = TestBed.createComponent(SeeAccommodationComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
