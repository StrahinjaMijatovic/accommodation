import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';  // Dodaj Router

@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.css']
})
export class RegisterComponent {
  registerForm: FormGroup;

  constructor(private fb: FormBuilder, private http: HttpClient, private router: Router) {  // Dodaj Router u constructor
    this.registerForm = this.fb.group({
      firstName: ['', [Validators.required, Validators.minLength(2)]],  
      lastName: ['', [Validators.required, Validators.minLength(2)]],   
      username: ['', [Validators.required, Validators.minLength(5)]],   
      password: ['', [Validators.required, Validators.minLength(8)]],  
      confirmPassword: ['', Validators.required],
      email: ['', [Validators.required, Validators.email]],  
      age: ['', [Validators.required, Validators.min(18), Validators.max(100)]],  
      country: ['', Validators.required],  
      role: ['', Validators.required]  
    }, { validator: this.passwordMatchValidator });
  }


  passwordMatchValidator(form: FormGroup) {
    return form.get('password')?.value === form.get('confirmPassword')?.value
      ? null : { mismatch: true };
  }

 
  get f() { return this.registerForm.controls; }

  onSubmit() {
    if (this.registerForm.valid) {
      this.http.post('http://localhost:8000/register', this.registerForm.value)
        .subscribe(
          response => {
            console.log('Registration successful', response);
            this.router.navigate(['/login']);  
          },
          error => {
            console.error('Registration error', error);
          }
        );
    } else {
      console.log('Form is invalid');
    }
  }
}
