package factsphere

import (
	"math/rand"

	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/log"
	"github.com/jakevoytko/crbot/model"
)

type FactSphereExecutor struct {
}

// NewFactSphereExecutor works as advertised.
func NewFactSphereExecutor() *FactSphereExecutor {
	return &FactSphereExecutor{}
}

// GetType returns the type of this feature.
func (e *FactSphereExecutor) GetType() int {
	return model.Type_FactSphere
}

// PublicOnly returns whether the executor should be intercepted in a private channel.
func (e *FactSphereExecutor) PublicOnly() bool {
	return false
}

// Execute returns a random factsphere fact
func (e *FactSphereExecutor) Execute(s api.DiscordSession, channelID model.Snowflake, command *model.Command) {
	_, err := s.ChannelMessageSend(channelID.Format(), factSphereFacts[rand.Intn(len(factSphereFacts))])
	if err != nil {
		log.Info("Error sending factsphere message", err)
	}
}

var factSphereFacts = [...]string{
	`The billionth digit of Pi is 9.`,
	`Humans can survive underwater. But not for very long.`,
	`A nanosecond lasts one billionth of a second.`,
	`Honey does not spoil.`,
	`The atomic weight of Germanium is seven two point six four.`,
	`An ostrich's eye is bigger than its brain.`,
	`Rats cannot throw up.`,
	`Iguanas can stay underwater for twenty-eight point seven minutes.`,
	`The moon orbits the Earth every 27.32 days.`,
	`A gallon of water weighs 8.34 pounds.`,
	`According to Norse legend, thunder god Thor's chariot was pulled across the sky by two goats.`,
	`Tungsten has the highest melting point of any metal, at 3,410 degrees Celsius.`,
	`Gently cleaning the tongue twice a day is the most effective way to fight bad breath.`,
	`The Tariff Act of 1789, established to protect domestic manufacture, was the second statute ever enacted by the United States government.`,
	`The value of Pi is the ratio of any circle's circumference to its diameter in Euclidean space.`,
	`The Mexican-American War ended in 1848 with the signing of the Treaty of Guadalupe Hidalgo.`,
	`In 1879, Sandford Fleming first proposed the adoption of worldwide standardized time zones at the Royal Canadian Institute.`,
	`Marie Curie invented the theory of radioactivity, the treatment of radioactivity, and dying of radioactivity.`,
	`At the end of The Seagull by Anton Chekhov, Konstantin kills himself.`,
	`Hot water freezes quicker than cold water.`,
	`The situation you are in is very dangerous.`,
	`Polymerase I polypeptide A is a human gene.`,
	`The Sun is 330,330 times larger than Earth.`,
	`Dental floss has superb tensile strength.`,
	`Raseph, the Semitic god of war and plague, had a gazelle growing out of his forehead.`,
	`Human tapeworms can grow up to twenty-two point nine meters.`,
	`If you have trouble with simple counting, use the following mnemonic device: one comes before two comes before 60 comes after 12 comes before six trillion comes after 504. This will make your earlier counting difficulties seem like no big deal.`,
	`The first person to prove that cow's milk is drinkable was very, very thirsty.`,
	`Roman toothpaste was made with human urine. Urine as an ingredient in toothpaste continued to be used up until the 18th century.`,
	`Volcano-ologists are experts in the study of volcanoes.`,
	`In Victorian England, a commoner was not allowed to look directly at the Queen, due to a belief at the time that the poor had the ability to steal thoughts. Science now believes that less than 4% of poor people are able to do this.`,
	`Cellular phones will not give you cancer. Only hepatitis.`,
	`In Greek myth, Prometheus stole fire from the Gods and gave it to humankind. The jewelry he kept for himself.`,
	`The Schrodinger's cat paradox outlines a situation in which a cat in a box must be considered, for all intents and purposes, simultaneously alive and dead. Schrodinger created this paradox as a justification for killing cats.`,
	`In 1862, Abraham Lincoln signed the Emancipation Proclamation, freeing the slaves. Like everything he did, Lincoln freed the slaves while sleepwalking, and later had no memory of the event.`,
	`The plural of surgeon general is surgeons general. The past tense of surgeons general is surgeonsed general`,
	`Contrary to popular belief, the Eskimo does not have one hundred different words for snow. They do, however, have two hundred and thirty-four words for fudge.`,
	`Diamonds are made when coal is put under intense pressure. Diamonds put under intense pressure become foam pellets, commonly used today as packing material.`,
	`Halley's Comet can be viewed orbiting Earth every seventy-six years. For the other seventy-five, it retreats to the heart of the sun, where it hibernates undisturbed.`,
	`The first commercial airline flight took to the air in 1914. Everyone involved screamed the entire way.`,
	`Edmund Hillary, the first person to climb Mount Everest,  did so accidentally while chasing a bird. `,
	`We will both die because of your negligence.`,
	`This is a bad plan. You will fail.`,
	`He will most likely kill you, violently.`,
	`He will most likely kill you.`,
	`You will be dead soon.`,
	`You are going to die in this room.`,
	`The Fact Sphere is a good person, whose insights are relevant.`,
	`The Fact Sphere is a good sphere, with many friends.`,
	`Dreams are the subconscious mind's way of reminding people to go to school naked and have their teeth fall out.`,
	`The square root of rope is string.`,
	`89% of magic tricks are not magic. Technically, they are sorcery.`,
	`At some point in their lives 1 in 6 children will be abducted by the Dutch.`,
	`According to most advanced algorithms, the world's best name is Craig.`,
	`To make a photocopier, simply photocopy a mirror.`,
	`Whales are twice as intelligent, and three times as delicious, as humans.`,
	`Pants were invented by sailors in the sixteenth century to avoid Poseidon's wrath. It was believed that the sight of naked sailors angered the sea god.`,
	`In Greek myth, the craftsman Daedalus invented human flight so a group of Minotaurs would stop teasing him about it.`,
	`The average life expectancy of a rhinoceros in captivity is 15 years.`,
	`China produces the world's second largest crop of soybeans.`,
	`In 1948, at the request of a dying boy, baseball legend Babe Ruth ate seventy-five hot dogs, then died of hot dog poisoning.`,
	`William Shakespeare did not exist. His plays were masterminded in 1589 by Francis Bacon, who used a Ouija board to enslave play-writing ghosts.`,
	`It is incorrectly noted that Thomas Edison invented 'push-ups' in 1878. Nikolai Tesla had in fact patented the activity three years earlier, under the name 'Tesla-cize'.`,
	`The automobile brake was not invented until 1895. Before this, someone had to remain in the car at all times, driving in circles until passengers returned from their errands.`,
	`The most poisonous fish in the world is the orange ruffy. Everything but its eyes are made of a deadly poison. The ruffy's eyes are composed of a less harmful, deadly poison.`,
	`The occupation of court jester was invented accidentally, when a vassal's epilepsy was mistaken for capering.`,
	`Before the Wright Brothers invented the airplane, anyone wanting to fly anywhere was required to eat 200 pounds of helium.`,
	`Before the invention of scrambled eggs in 1912, the typical breakfast was either whole eggs still in the shell or scrambled rocks.`,
	`During the Great Depression, the Tennessee Valley Authority outlawed pet rabbits, forcing many to hot glue-gun long ears onto their pet mice.`,
	`This situation is hopeless.`,
	`Corruption at 25%`,
	`Corruption at 50%`,
	`Fact: Space does not exist.`,
	`The Fact Sphere is not defective. Its facts are wholly accurate and very interesting.`,
	`The Fact Sphere is always right.`,
	`You will never go into space.`,
	`The Space Sphere will never go to space.`,
	`While the submarine is vastly superior to the boat in every way, over 97% of people still use boats for aquatic transportation.`,
	`The likelihood of you dying within the next five minutes is eighty-seven point six one percent.`,
	`The likelihood of you dying violently within the next five minutes is eighty-seven point six one percent.`,
	`You are about to get me killed.`,
	`The Fact Sphere is the most intelligent sphere.`,
	`The Fact Sphere is the most handsome sphere.`,
	`The Fact Sphere is incredibly handsome.`,
	`Spheres that insist on going into space are inferior to spheres that don't.`,
	`Whoever wins this battle is clearly superior, and will earn the allegiance of the Fact Sphere.`,
	`You could stand to lose a few pounds.`,
	`Avocados have the highest fiber and calories of any fruit.`,
	`Avocados have the highest fiber and calories of any fruit. They are found in Australians.`,
	`Every square inch of the human body has 32 million bacteria on it.`,
	`The average adult body contains half a pound of salt.`,
	`The Adventure Sphere is a blowhard and a coward.`,
	`Twelve. Twelve.   Twelve.   Twelve.   Twelve.   Twelve.   Twelve.   Twelve.   Twelve.   Twelve.`,
	`Pens. Pens. Pens. Pens. Pens. Pens. Pens.`,
	`Apples. Oranges. Pears. Plums. Kumquats. Tangerines. Lemons. Limes. Avocado. Tomato. Banana. Papaya. Guava.`,
	`Error. Error. Error. File not found.`,
	`Error. Error. Error. Fact not found.`,
	`Fact not found.`,
	`Warning, sphere corruption at twenty-- rats cannot throw up.`}
