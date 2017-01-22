import { Component, OnInit, ViewChild } from '@angular/core';
import { NgForm } from '@angular/forms';
import { Router, ActivatedRoute, Params } from '@angular/router';

import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { Observable } from 'rxjs/Observable';

import { NotificationsService } from 'angular2-notifications';

import { CreateResponse } from '../api/create-response';
import { UpdateResponse } from '../api/update-response';

import { Member } from "./member";
import { MemberService } from './member.service';

@Component({
    templateUrl: 'app/members/member.component.html',
    styleUrls: ['app/members/member.component.css']
})
export class MemberComponent {
    constructor(
        private router: Router,
        private route: ActivatedRoute,
        private memberService: MemberService,
        private notificationService: NotificationsService
    ) { }

    @ViewChild('memberForm') memberForm: NgForm;

    member: Member = new Member();
    errors: string[] = [];

    ngOnInit(): void {
        this.route.params.forEach((params: Params) => {
            if (params['id'] == 'new') {
                this.member = new Member;
                return;
            }

            let id = +params['id'];
            if (this.member.id != id) {
                // reset current member to avoid accidents during loading of new member
                this.member = new Member;
                // load requested member
                this.memberService.getMember(id)
                    .then((member: Member) => {
                        this.member = member;
                        this.memberForm.control.markAsPristine();
                    })
                    .catch((error: any) => alert(error));
            }
        });
    }

    private save() {
        this.errors = [];
        if (this.member.id == undefined) {
            this.memberService.postMember(this.member).then((res: CreateResponse) => {
                if (res.hasErrors()) {
                    this.handleSaveErrors(res.errors);
                    return;
                }
                this.router.navigate(['member', res.id]);
                this.notificationService.success('Nieuw lid opgeslagen', `'${this.member.initials} ${this.member.lastName}' is succesvol opgeslagen.`);
            }).catch((error: any) => alert(error));
        } else {
            this.memberService.putMember(this.member.id, this.member).then((res: UpdateResponse) => {
                if (res.hasErrors()) {
                    this.handleSaveErrors(res.errors);
                    return;
                }
                this.memberForm.control.markAsPristine();
                this.notificationService.success('Aanpassingen opgeslagen', `Aanpassingen aan '${this.member.initials} ${this.member.lastName}' zijn opgeslagen.`);
            }).catch((error: any) => alert(error));
        }
    }

    private handleSaveErrors(errors: string[]) {
        errors.forEach((error: any) => {
            switch (error) {
                case 'rutte:invalid_email_address':
                case 'pgerr:accounts_uq_email':
                case 'pgerr:check_violation:accounts_check_age_over_14':
                    // handled by form
                    break;
                case 'pgerr:datetime_field_overflow':
                    alert('Datum fout'); // TODO: dialog
                    break;
                default:
                    alert('unhandled error: ' + error);
                    break;
            }
        });
        this.errors = errors;
        return;
    }
}
