import { Component, OnInit } from '@angular/core';
import { Http, Response } from '@angular/http';
import { Router, ActivatedRoute, Params } from '@angular/router';

import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { Observable } from 'rxjs/Rx';
import { TypedJSON, JsonObject, JsonMember } from 'typedjson-npm';
import { SessionStorage } from 'ng2-webstorage';

import { ConfigService } from '../common/config.service';
import { Date } from '../common/date';
import { CountryList, Country } from '../common/countries';

class MemberRegistration {
	// step indicates the last saved step for this registration
	step: number = 1;

	// step1 form fields
	givenName: string;
	lastName: string;
	birthdate: Date;
	email: string;

	// step1 api call results
	id: number = null;
	registrationToken: string = null;

	// step2 form fields
	// streetname, city and province can be autocompleted by postalcode
	postalcode: string;
	housenumber: number;
	houstenumberSuffix: string;
	streetname: string;
	city: string;
	province: string;
	countryCode: string = null; // selected in pre-step2 with buttons. Must be set to null by default to diff between !'NL' and not chosen yet.
}

class PostalcodeHousenumber {
	postalcode: string;
	housenumber: number;
}

@JsonObject()
class ResponseStep1 {
	@JsonMember() id: number;
	@JsonMember() registrationToken: string;

	// errors array, the request was ok when length==0
	// contains non-human-friendly errorcodes
	@JsonMember({ elements: String }) errors: string[];
}

@JsonObject()
class ResponseStep2 {
	// fields returned when countryCode=='NL', using postcode API serverside
	@JsonMember() streetname: string;
	@JsonMember() city: string;
	@JsonMember() province: string;

	// errors array, the request was ok when length==0
	// contains non-human-friendly errorcodes
	@JsonMember({ elements: String }) errors: string[];
}

@JsonObject()
class ResponseStep3 {
	@JsonMember() payment_url: string;

	// errors array, the request was ok when length==0
	// contains non-human-friendly errorcodes
	@JsonMember({ elements: String }) errors: string[];
}

@Component({
	styleUrls: [`app/register/register.component.css`],
	templateUrl: `app/register/register.component.html`,
})
export class RegisterComponent {
	/**
	 * steps:
	 * 1: name, email address
	 * 2: address
	 * 3: validate + confirm {incasso, stemrecht}
	 * 4: redirect to mollie
	 * 5: mollie success or failure
	 **/
	step = 1;

	// saveInProgress indicates that a save call is being made. locks forms
	saveInProgress: boolean = null;

	// member details during registration
	@SessionStorage('register.member')
	public member: MemberRegistration; // TODO: needs to be public?

	// errors returned by the server during one of the steps.
	errors: string[] = [];

	// country list
	countries: Observable<Country[]>;

	// postalcode filterstream
	private postalcodeHousenumberUpdateStream = new BehaviorSubject<PostalcodeHousenumber>(null);
	postalcodeHousenumberUpdated() { this.postalcodeHousenumberUpdateStream.next({ postalcode: this.member.postalcode, housenumber: this.member.housenumber }); }

	constructor(
		private router: Router,
		private route: ActivatedRoute,
		private http: Http,
		private config: ConfigService
	) {
		if (this.member == null) {
			this.member = new MemberRegistration();
		}
		console.dir(this.member);

		this.countries = this.http.get(this.config.apiURL + '/api/country-list').map((res: Response) => TypedJSON.parse(res.text(), CountryList).countries);

		// setup postalcode filterstream
		this.postalcodeHousenumberUpdateStream
			.filter((ph: PostalcodeHousenumber) => ph != null)
			.filter((ph: PostalcodeHousenumber) => ph.postalcode.replace(`^([0-9]{4}) ?([a-zA-Z]{2})$`, `$1$2`).length == 6)
			.filter((ph: PostalcodeHousenumber) => +ph.housenumber > 0)
			.debounceTime(400)
			.distinctUntilChanged()
			.forEach(() => this.saveStep2Partial());
	}

	ngOnInit(): void {

		this.route.params.forEach((params: Params) => {
			if (!params['step']) {
				this.step = 1;
				return;
			}

			this.step = +params['step'];
		});
	}

	private saveStep1(empForm: any, event: Event) {
		event.preventDefault();
		this.errors = [];

		// skip post if step has already been processed
		if (this.member.step > 1) {
			this.router.navigate(['lid-worden', 2]);
			return;
		}

		// post step1
		this.saveInProgress = true;
		this.http.post(this.config.apiURL + `/api/register/step1`, this.member).toPromise().then((res: Response) => {
			let responseStep1 = TypedJSON.parse(res.text(), ResponseStep1);
			if (responseStep1.errors.length > 0) {
				this.handleSaveErrorsStep1(responseStep1.errors);
				this.saveInProgress = null;
				return;
			}
			this.member.id = responseStep1.id;
			this.member.registrationToken = responseStep1.registrationToken;
			this.member.step = 2;
			this.member = this.member; // trigger WebStorage decorator save

			// unlock forms, go to next step.
			this.saveInProgress = null;
			this.router.navigate(['lid-worden', 2]);
		}).catch((error: any) => alert(error));
	}

	private handleSaveErrorsStep1(errors: string[]) {
		errors.forEach((error: any) => {
			switch (error) {
				case 'rutte:invalid_email_address':
				case 'pgerr:accounts_uq_email':
				case 'pgerr:check_violation:accounts_check_age_over_14':
					// handled by form
					break;
				// case 'pgerr:datetime_field_overflow':
				// 	alert('Datum fout'); // TODO: dialog
				// 	break;
				default:
					alert('unhandled error: ' + error); // TODO: dialog
					break;
			}
		});
		this.errors = errors;
		return;
	}

	private saveStep2Partial(): Promise<boolean> {
		this.errors = [];
		return this.http.post(this.config.apiURL + `/api/register/step2`, this.member).toPromise().then((res: Response) => {
			let responseStep2 = TypedJSON.parse(res.text(), ResponseStep2);
			if (responseStep2.errors.length > 0) {
				this.handleSaveErrorsStep2(responseStep2.errors);
				return false;
			}
			if (this.member.countryCode = 'NL') {
				this.member.streetname = responseStep2.streetname;
				this.member.city = responseStep2.city;
				this.member.province = responseStep2.province;
			}
			return true
		}).catch((error: any) => alert(error));
	}

	private saveStep2(empForm: any, event: Event) {
		event.preventDefault();

		// skip post if step has already been processed
		if (this.member.step > 2) {
			this.router.navigate(['lid-worden', 3]);
			return;
		}

		// post step1
		this.saveInProgress = true;
		this.saveStep2Partial().then((success: boolean) => {
			if (!success) {
				this.saveInProgress = null;
				return;
			}
			this.member.step = 3;
			this.member = this.member; // trigger WebStorage decorator save

			// unlock forms, go to next step.
			this.saveInProgress = null;
			this.router.navigate(['lid-worden', 3]);
		}).catch((error: any) => alert(error));
	}

	private handleSaveErrorsStep2(errors: string[]) {
		errors.forEach((error: any) => {
			switch (error) {
				case 'rutte:PostcodeNl_Service_PostcodeAddress_AddressNotFoundException':
					// handled by form
					break;
				default:
					alert('unhandled error: ' + error); // TODO: dialog
					break;
			}
		});
		this.errors = errors;
		return;
	}

	private saveStep3(empForm: any, event: Event) {
		event.preventDefault();

		this.http.post(this.config.apiURL + `/api/register/step3`, this.member).toPromise().then((res: Response) => {
			let responseStep3 = TypedJSON.parse(res.text(), ResponseStep3);
			if (responseStep3.errors.length > 0) {
				this.handleSaveErrorsStep2(responseStep3.errors);
			}
			console.log(responseStep3.payment_url);
		}).catch((error: any) => alert(error));
	}
}
