import { Component } from "@angular/core";

import { Auth } from "./auth.service";

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

		<div *ngIf="!auth.isLoggedIn()" >
			<input [(ngModel)]="id" type="number" >
			<button (click)="login()" >Log in</button>
		</div>
	`,
})
export class LoginComponent {
	public id: number;

	// feedback messages
	public successMessage: string;
	public errorMessage: string;

	constructor(public auth: Auth) { }

	// login tries to authenticate using the current form values
	public login() {
		this.successMessage = ``;
		this.errorMessage = ``;
		this.auth.tryLogin(this.id).subscribe(
			() => this.successMessage = "all good!",
			(error: any) => this.errorMessage = error
		);
	}
}
