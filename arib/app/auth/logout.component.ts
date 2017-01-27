import { Component } from '@angular/core';

import { Auth } from './auth.service';

@Component({
	template: `
		{{errorMessage}} | {{successMessage}}
	`,
})
export class LogoutComponent {
	id: number;

	// feedback messages
	successMessage: string;
	errorMessage: string;

	constructor(private auth: Auth) {
		this.logout()
	}

	// login tries to authenticate using the current form values
	logout() {
		this.successMessage = ``;
		this.errorMessage = ``;
		this.auth.logout().subscribe(
			(data: any) => {
				if (data.error) {
					switch (data.error) {
						case 'rutte:no_session_found':
							this.errorMessage = 'Er was geen bestaande inlogsessie. Mogelijk was u al uitgelogd';
							break;
						default:
							alert('unhandled error ' + data.error);
							break;
					}
					return
				}
				this.successMessage = 'U bent uitgelogd.'
			},
			(error: any) => this.errorMessage = error
		);
	}
}
