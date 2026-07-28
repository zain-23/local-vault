import { faker } from "@faker-js/faker";
import type { Member, MemberRole } from "#/features/members/types";

// Deterministic seed so the mock roster is stable across reloads (no layout jitter).
faker.seed(42);

function makeMember(role: MemberRole): Member {
	const name = faker.person.fullName();
	return {
		userId: faker.string.uuid(),
		name,
		email: faker.internet
			.email({ firstName: name.split(" ")[0] })
			.toLowerCase(),
		avatarUrl: faker.image.avatarGitHub(),
		role,
		joinedAt: faker.date.past({ years: 2 }).toISOString(),
	};
}

// One owner, a couple of admins, the rest members — mirrors a real small team.
const ROLES: MemberRole[] = [
	"owner",
	"admin",
	"admin",
	"member",
	"member",
	"member",
	"member",
	"member",
	"member",
	"member",
];

export const MOCK_MEMBERS: Member[] = ROLES.map(makeMember);
