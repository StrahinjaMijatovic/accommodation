import { ComponentFixture, TestBed } from '@angular/core/testing';

import { HostNotificationsComponent } from './host-notifications.component';

describe('HostNotificationsComponent', () => {
  let component: HostNotificationsComponent;
  let fixture: ComponentFixture<HostNotificationsComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [HostNotificationsComponent]
    });
    fixture = TestBed.createComponent(HostNotificationsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
