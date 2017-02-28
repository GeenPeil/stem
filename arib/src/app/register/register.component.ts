import { Component, OnInit } from "@angular/core";
import { Http, Response } from "@angular/http";
import { Router, ActivatedRoute, Params } from "@angular/router";

import { BehaviorSubject } from "rxjs/BehaviorSubject";
import { Observable } from "rxjs/Rx";
import { TypedJSON, JsonObject, JsonMember } from "typedjson-npm";
import { SessionStorage } from "ng2-webstorage";

import { ConfigService } from "../common/config.service";
import { Date } from "../common/date";
import { CountryList, Country } from "../common/countries";

class MemberRegistration {
	// step indicates the last saved step for this registration
	public step: number = 1;

	// step1 form fields
	public givenName: string;
	public lastName: string;
	public birthdate: Date;
	public email: string;
	public phonenumber: string;

	// step1 api call results
	public id: number = null;
	public registrationToken: string = null;

	// step2 form fields
	// streetname, city and province can be autocompleted by postalcode
	public postalcode: string;
	public housenumber: number;
	public houstenumberSuffix: string;
	public streetname: string;
	public city: string;
	public province: string;
	// countryCode is selected in pre-step2 with buttons. It must be set to null by default to diff between !"NL" and not chosen yet.
	public countryCode: string = null;

	public kalenderjaar: boolean = false;
	public stemrecht: boolean = false;
}

class PostalcodeHousenumber {
	public postalcode: string;
	public housenumber: number;
}

@JsonObject()
class ResponseStep1 {
	@JsonMember() public id: number;
	@JsonMember() public registrationToken: string;

	// errors array, the request was ok when length==0
	// contains non-human-friendly errorcodes
	@JsonMember({ elements: String }) public errors: string[];
}

@JsonObject()
class ResponseStep2 {
	// fields returned when countryCode==="NL", using postcode API serverside
	@JsonMember() public streetname: string;
	@JsonMember() public city: string;
	@JsonMember() public province: string;

	// errors array, the request was ok when length==0
	// contains non-human-friendly errorcodes
	@JsonMember({ elements: String }) public errors: string[];
}

@JsonObject()
class ResponseStep3 {
	@JsonMember() public paymentURL: string;

	// errors array, the request was ok when length==0
	// contains non-human-friendly errorcodes
	@JsonMember({ elements: String }) public errors: string[];
}

@JsonObject()
class ResponseCheckPayment {
	@JsonMember() public paid: boolean;

	// errors array, the request was ok when length==0
	// contains non-human-friendly errorcodes
	@JsonMember({ elements: String }) public errors: string[];
}

@Component({
	styleUrls: [`register.component.css`],
	templateUrl: `register.component.html`,
})
export class RegisterComponent implements OnInit {
	/**
	 * steps:
	 * 1: name, email address
	 * 2: address
	 * 3: validate + confirm {incasso, stemrecht}
	 * 4: redirect to mollie
	 * 5: mollie success or failure
	 */
	public step = 1;

	// saveInProgress indicates that a save call is being made. locks forms
	public saveInProgress: boolean = null;

	// viaLink
	public viaLink: boolean = false;

	// member details during registration
	@SessionStorage("register.member")
	public member: MemberRegistration; // TODO: needs to be public?

	// errors returned by the server during one of the steps.
	public errors: string[] = [];

	// country list
	public countries: Observable<Country[]>;

	// checkPaymentToken is set when returning from a payment.
	public checkPaymentToken: string = null;
	public checkPaymentSuccess: boolean = null;

	// postalcode filterstream
	private postalcodeHousenumberUpdateStream = new BehaviorSubject<PostalcodeHousenumber>(null);

	constructor(
		private router: Router,
		private route: ActivatedRoute,
		private http: Http,
		private config: ConfigService
	) {
		if (this.member === null) {
			this.member = new MemberRegistration();
			let params = this.parseQueryString();
			console.log(params["gn"]);
			if (params["gn"] !== undefined) {
				this.member.givenName = params["gn"];
				if (params["ln"] !== undefined) {
					this.member.lastName = params["ln"];
				}
				this.viaLink = true;
				this.member = this.member;
			}
		}
		console.dir(this.member);

		this.countries = this.http.get(this.config.apiURL + "/api/country-list").map(
			(res: Response) => TypedJSON.parse(res.text(), CountryList).countries
		);

		// setup postalcode filterstream
		this.postalcodeHousenumberUpdateStream
			.filter((ph: PostalcodeHousenumber) => ph !== null)
			.filter((ph: PostalcodeHousenumber) => ph.postalcode.replace(`^([0-9]{4}) ?([a-zA-Z]{2})$`, `$1$2`).length === 6)
			.filter((ph: PostalcodeHousenumber) => +ph.housenumber > 0)
			.debounceTime(400)
			.distinctUntilChanged()
			.forEach(() => this.saveStep2Partial());
	}

	public ngOnInit(): void {
		this.route.params.forEach((params: Params) => {
			if (!params["step"]) {
				this.step = 1;
				return;
			}

			if (params["step"] === "check-payment") {
				this.step = 3;
				this.checkPaymentToken = this.parseQueryString()["token"];
				this.checkPayment();
				return;
			}

			this.step = +params["step"];
		});
	}

	public postalcodeHousenumberUpdated() {
		this.postalcodeHousenumberUpdateStream.next({ postalcode: this.member.postalcode, housenumber: this.member.housenumber });
	}

	public saveStep1(empForm: any, event: Event) {
		event.preventDefault();
		this.errors = [];

		// skip post if step has already been processed
		if (this.member.step > 1) {
			this.router.navigate(["lid-worden", 2]);
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
			this.router.navigate(["lid-worden", 2]);
		}).catch((error: any) => alert(error));
	}

	private handleSaveErrorsStep1(errors: string[]) {
		errors.forEach((error: any) => {
			switch (error) {
				case "rutte:invalid_email_address":
				case "pgerr:check_violation:accounts_check_age_over_14":
					// handled by form
					break;
				// case "pgerr:datetime_field_overflow":
				// 	alert("Datum fout"); // TODO: dialog
				// 	break;
				default:
					alert("unhandled error: " + error); // TODO: dialog
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
			if (this.member.countryCode === "NL") {
				this.member.streetname = responseStep2.streetname;
				this.member.city = responseStep2.city;
				this.member.province = responseStep2.province;
			}
			return true;
		}).catch((error: any) => alert(error));
	}

	public saveStep2(empForm: any, event: Event) {
		event.preventDefault();

		// skip post if step has already been processed
		if (this.member.step > 2) {
			this.router.navigate(["lid-worden", 3]);
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
			this.router.navigate(["lid-worden", 3]);
		}).catch((error: any) => alert(error));
	}

	private handleSaveErrorsStep2(errors: string[]) {
		errors.forEach((error: any) => {
			switch (error) {
				case "rutte:PostcodeNl_Service_PostcodeAddress_AddressNotFoundException":
					// handled by form
					break;
				default:
					alert("unhandled error: " + error); // TODO: dialog
					break;
			}
		});
		this.errors = errors;
		return;
	}

	public saveStep3(empForm: any, event: Event) {
		event.preventDefault();

		this.http.post(this.config.apiURL + `/api/register/step3`, this.member).toPromise().then((res: Response) => {
			let responseStep3 = TypedJSON.parse(res.text(), ResponseStep3);
			if (responseStep3.errors.length > 0) {
				this.handleSaveErrorsStep3(responseStep3.errors);
				return;
			}
			this.member = this.member; // trigger WebStorage decorator save
			window.location.assign(responseStep3.payment_url);
		}).catch((error: any) => alert(error));
	}

	private handleSaveErrorsStep3(errors: string[]) {
		errors.forEach((error: any) => {
			switch (error) {
				case "pgerr:not_null_violation:account_id":
					this.member = null;
					alert("De registratie is niet succesvol verlopen. Probeer het a.u.b. opnieuw.");
					this.router.navigate(["lid-worden", 1]);
					break;
				default:
					alert("unhandled error: " + error); // TODO: dialog
					break;
			}
		});
		this.errors = errors;
		return;
	}

	private checkPayment() {
		this.http.get(this.config.apiURL + `/api/register/check-payment?token=` + this.checkPaymentToken).toPromise().then((res: Response) => {
			let responseCheckPayment = TypedJSON.parse(res.text(), ResponseCheckPayment);
			if (responseCheckPayment.errors.length > 0) {
				this.handleErrorsCheckPayment(responseCheckPayment.errors);
				return;
			}
			this.checkPaymentSuccess = responseCheckPayment.paid;
			if (this.checkPaymentSuccess) {
				this.member = null;
			}
		}).catch((error: any) => alert(error));
	}

	private handleErrorsCheckPayment(errors: string[]) {
		errors.forEach((error: any) => {
			switch (error) {
				case "rutte:invalid_payment_token":
					alert("Ongeldige URL");
					this.router.navigate(["lid-worden", 1]);
					break;
				default:
					alert("unhandled error: " + error); // TODO: dialog
					break;
			}
		});
		this.errors = errors;
	}

	private parseQueryString() {
		let str = window.location.search;
		let objURL = {};

		str.replace(
			new RegExp("([^?=&]+)(=([^&]*))?", "g"),
			($0, $1, $2, $3): string => {
				objURL[$1] = $3;
				return "";
			}
		);
		return objURL;
	};
}
