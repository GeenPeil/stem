import { Component } from '@angular/core';

import { Auth } from './auth.service';

@Component({
	template: `
		<div *ngIf="errorMessage" >
			Er was een probleem.<br/>
			{{errorMessage}}
		</div>

		<div *ngIf="successMessage" >
			{{successMessage}}<br/>
			<a routerLink="/profile" >Ga naar je profiel.</a>
		</div>

		<div *ngIf="!auth.IsLoggedIn()" >
			<input [(ngModel)]="id" type="number" >
			<button (click)="login()" >Log in</button>
		</div>
	`,
})
export class LoginComponent {
	id: number;

	// feedback messages
	successMessage: string;
	errorMessage: string;

	constructor(private auth: Auth) { }

	// login tries to authenticate using the current form values
	login() {
		this.successMessage = ``;
		this.errorMessage = ``;
		this.auth.tryLogin(this.id).subscribe(
			() => this.successMessage = 'all good!',
			(error: any) => this.errorMessage = error
		);
	}
}
