import { Component } from '@angular/core';

import { Api } from '../api/api.service';

import { Profile } from './profile';

@Component({
	template: `
		<div *ngIf="error">
			<h1>Could not load your profile</h1>
			<p>error: {{error}}</p>
			<button (click)="fetchProfile()">Retry</button>
		</div>
		<div *ngIf="profile">
			<h1>{{profile.nickname}} (member ID: {{profile.id}})</h1>
			<p>Email: {{profile.email}}</p>
		</div>
	`,
})
export class ProfileComponent {
	error: Error;
	profile: Profile;

	constructor(private api: Api) {
		this.fetchProfile();
	}

	private fetchProfile() {
		this.api.get(Profile, "private/profile").subscribe(
			(profile: Profile) => this.profile = profile,
			(error: Error) => this.error = error
		);
	}
}
