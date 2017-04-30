-- database
--

create type author_type as enum ('admin', 'user');
create type author_status as enum ('active', 'inactive');

create table if not exists author (
  handle varchar not null unique,
  email varchar not null unique,
  password varchar not null,
  status author_status not null default 'active',
  type author_type not null default 'user'
);

insert into author (email, handle, password)
  values ('keith@example.com', 'keith', 'test1234');
insert into author (email, handle, password)
  values ('xan@example.com', 'xan', 'test1234');

create type post_status as enum ('published', 'draft');

create table if not exists post (
  id serial primary key,
  uuid varchar not null unique,
  author varchar references author(handle) not null,
  date_created timestamp with time zone default current_timestamp,
  date_updated timestamp with time zone default current_timestamp,
  status post_status not null default 'draft',
  slugline text not null unique,
  text text default ''
);

insert into post (author, uuid, slugline, text) values ('keith',
'b8ec300c-b0f1-4338-b2a7-ca5f06c1fe33',
'Xena of Amphipolous',
'#Xena of Amphipolous

Xena is a fictional
character from Robert Tapert''s Xena: Warrior Princess
franchise. Co-created by Tapert and John Schulian, she first appeared
in the 1995–1999 television series Hercules: The Legendary Journeys,
before going on to appear in Xena: Warrior Princess TV show and
subsequent comic book of the same name. The Warrior Princess has also
appeared in the spin-off animated movie The Battle for Mount Olympus,
as well as numerous non-canon expanded universe material, such as
books and video games. Xena was played by New Zealand actress Lucy
Lawless.

Xena is the protagonist of the story, and the series depicts her on a
quest to redeem herself for her dark past by using her formidable
fighting skills to help people. Xena was raised as the daughter of
Cyrene and Atrius in Amphipolis, though the episode "The Furies"
raised the possibility that Ares might be Xena''s biological father,
but it is never pursued further. She had two brothers, the younger of
whom is dead; she visits his grave to speak with him in "Sins of the
Past." In Hercules: The Legendary Journeys, during her two first
episodes, Xena was a villain, but in the third episode she appears in,
she joins Hercules to defeat Darphus, who had taken her army. Aware
that the character of Xena had been very successful with the public in
the three Hercules: The Legendary Journeys episodes, the producers of
the series decided to create a spin-off series based on her
adventures. Later in Xena: Warrior Princess she is joined by
Gabrielle, a small town bard. Together they go up against ruthless
warlords and gods in the ancient mythological world.

> Xena has been credited by many, including Buffy the Vampire Slayer
  creator Joss Whedon, with blazing the trail for a new generation of
  female action heroes such as Buffy, Max of Dark Angel, Sydney
  Bristow of Alias, and Beatrix Kiddo a.k.a. the Bride in Quentin
  Tarantino''s Kill Bill. The director Quentin Tarantino is also a
  fan of Xena. It is interesting to note that after serving as Lucy
  Lawless'' stunt double on Xena, stunt woman Zoë E. Bell was recruited
  to be Uma Thurman''s stunt double in Tarantino''s Kill Bill. By
  helping to pave the way for female action heroes in television and
  film, "Xena" also strengthened the stunt woman profession. David
  Eick, one of the co-developers of the Xena series, was also the
  executive producer of Battlestar Galactica,] which also features
  strong female characters, and Lucy Lawless in a recurring role.

The character Gabrielle, introduced in the first episode, becomes
Xena''s greatest ally, best friend, and soulmate. Gabrielle came from a
small village in Greece called Potidaea. She craved to escape from the
boring and dull village. She latched onto Xena in episode 1 as a way
of leaving the village, to travel on adventures. Her initial naïveté
for the first 3 seasons and her talkative nature helped to balance
Xena''s pessimistic mentality. While Xena''s character to an extent
alters subtly through the series, Gabrielle''s character goes through
substantial development and change especially in seasons 3 and
4. Through their friendship/relationship Xena recognizes the value of
the "greater good" and the sacrifices that must be made to accomplish
it (a central theme in the series in later seasons).  ' );

insert into
       post (author, uuid, slugline, text) values ('xan',
'd08e9b3e-466b-423f-bd41-761f984bd0d0',
'Zeno of Elea',
'# Zeno of Elia

Little is known for certain about Zeno''s life. Although written
nearly a century after Zeno''s death, the primary source of
biographical information about Zeno is Plato''s Parmenides and he is
also mentioned in Aristotle''s Physics. In the dialogue of
Parmenides, Plato describes a visit to Athens by Zeno and Parmenides,
at a time when Parmenides is "about 65," Zeno is "nearly 40" and
Socrates is "a very young man". Assuming an age for Socrates of
around 20, and taking the date of Socrates'' birth as 469 BC gives an
approximate date of birth for Zeno of 490 BC. Plato says that Zeno was
"tall and fair to look upon" and was "in the days of his youth ...
reported to have been beloved by Parmenides."

Other perhaps less reliable details of Zeno''s life are given by
Diogenes Laërtius in his Lives and Opinions of Eminent
Philosophers,[6] where it is reported that he was the son of
Teleutagoras, but the adopted son of Parmenides, was "skilled to argue
both sides of any question, the universal critic," and that he was
arrested and perhaps killed at the hands of a tyrant of Elea.

- This is just to say
- That I''ve added a list
- In this document to
- See what it looks like when
- Interpreted

According to Diogenes Laertius, Zeno conspired to overthrow Nearchus
the tyrant.[9] Eventually, Zeno was arrested and tortured.
According to Valerius Maximus, when he was tortured to reveal the name
of his colleagues in conspiracy, Zeno refused to reveal their names,
although he said he did have a secret that would be advantageous for
Nearchus to hear. When Nearchus leaned in to listen to the secret,
Zeno bit his ear. He "did not let go until he lost his life and the
tyrant lost that part of his body." Within Men of the Same
Name, Demetrius said it was the nose that was bit off instead.

Zeno may have also interacted with other tyrants. According to
Laertius, Heraclides Lembus, within his Satyrus, said these events
occurred against Diomedon instead of Nearchus. Valerius Maximus
recounts a conspiracy against the tyrant Phalaris, but this would be
impossible as Phalaris had died before Zeno was even born.
According to Plutarch, Zeno attempted to kill the tyrant
Demylus. After failing, he had, "with his own teeth bit off his
tongue, he spit it in the tyrant’s face."
');

-- To view all the enums
create or replace view vw_enums as
select
  t.typname, e.enumlabel, e.enumsortorder
from
  pg_enum e
join
  pg_type t ON e.enumtypid = t.oid;
