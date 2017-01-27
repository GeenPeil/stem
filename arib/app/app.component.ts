import { Component } from '@angular/core';

@Component({
	selector: 'gp-arib',
	template: `
		<template ngbModalContainer></template>

		<nav class="navbar navbar-toggleable navbar-light bg-faded">
			<button class="navbar-toggler navbar-toggler-right" type="button" (click)="isNavbarCollapsed = !isNavbarCollapsed"
					aria-controls="navbarNavAltMarkup" aria-expanded="false" aria-label="Toon menu">
				<span class="navbar-toggler-icon"></span>
			</button>
			<a class="navbar-brand" href="#">GeenPeil</a>
			<div [ngbCollapse]="isNavbarCollapsed" class="collapse navbar-collapse" id="navbarNavAltMarkup">
				<div class="navbar-nav">
					<a class="nav-item nav-link" routerLinkActive="active" routerLink="/does-not-exist">Does not exist</a>
					<a class="nav-item nav-link" routerLinkActive="active" routerLink="/profile">Profile</a>
					<a class="nav-item nav-link" routerLinkActive="active" routerLink="/login">Login</a>
				</div>
			</div>
		</nav>

		<router-outlet></router-outlet>
	`,
})
export class AppComponent {
	isNavbarCollapsed = true;
}
