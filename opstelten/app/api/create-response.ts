import { JsonObject, JsonMember } from 'typedjson-npm';

@JsonObject()
export class CreateResponse {
    @JsonMember id: number;
    @JsonMember({ elements: String }) errors: string[];

    hasErrors(): boolean {
        return this.errors.length > 0;
    }
}
