package main

import (
	"errors"
	"fmt"
	"sort"
)

type card int

const (
	BlinkCard card = iota
	DigCard
	TeleportCard
	TeleportOtherCard
	HealWoundsCard
	MagicCard
	DescentCard
	SwiftnessCard
	SwappingCard
	ShadowsCard
	FogCard
	MagicMappingCard
	SensingCard
	WallsCard
	SlowingCard
	SleepingCard
	NoiseCard
	ObstructionCard
	FireCard
	ConfusionCard
	LignificationCard
)

func (g *game) DrawCard(n int) {
	cards := []card{
		BlinkCard,
		DigCard,
		TeleportCard,
		TeleportOtherCard,
		HealWoundsCard,
		MagicCard,
		DescentCard,
		SwiftnessCard,
		SwappingCard, // XXX: swap with random monster?
		ShadowsCard,
		FogCard,
		MagicMappingCard,
		SensingCard,
		WallsCard,
		SlowingCard,
		SleepingCard,
		NoiseCard,
		ObstructionCard,
		FireCard,
		ConfusionCard,
		LignificationCard,
	}
loop:
	for {
		c := cards[RandInt(len(cards))]
		for _, oc := range g.GeneratedCards {
			if c == oc && RandInt(5) > 0 {
				continue loop
			}
		}
		switch c {
		case MagicMappingCard, SensingCard, DescentCard:
			// rare cards appear less often
			if RandInt(2) == 0 {
				continue loop
			}
		}
		g.Hand[n] = c
		g.GeneratedCards = append(g.GeneratedCards, c)
		if len(g.GeneratedCards) > 4 {
			g.GeneratedCards = g.GeneratedCards[len(g.GeneratedCards)-4:]
		}
	}
}

func (g *game) UseCard(n int, ev event) (err error) {
	if g.Player.HasStatus(StatusNausea) {
		return errors.New("You cannot use cards while sick.")
	}
	c := g.Hand[n]
	if c.MPCost() > g.Player.MP {
		return errors.New("Not enough magic points for using this rod.")
	}
	switch c {
	case BlinkCard:
		err = g.EvokeBlink(ev)
	case TeleportCard:
		err = g.EvokeTeleport(ev)
	case DigCard:
		err = g.EvokeDig(ev)
	case TeleportOtherCard:
		err = g.EvokeTeleportOther(ev)
	case HealWoundsCard:
		err = g.EvokeHealWounds(ev)
	case MagicCard:
		err = g.EvokeRefillMagic(ev)
	case DescentCard:
		err = g.EvokeDescent(ev)
	case SwiftnessCard:
		err = g.EvokeSwiftness(ev)
	case SwappingCard:
		err = g.EvokeSwapping(ev)
	case ShadowsCard:
		err = g.EvokeShadows(ev)
	case FogCard:
		err = g.EvokeFog(ev)
	case MagicMappingCard:
		err = g.EvokeMagicMapping(ev)
	case SensingCard:
		err = g.EvokeSensing(ev)
	case WallsCard:
		err = g.EvokeWalls(ev)
	case SlowingCard:
		err = g.EvokeSlowing(ev)
	case SleepingCard:
		err = g.EvokeSleeping(ev)
	case NoiseCard:
		err = g.EvokeNoise(ev)
	case ConfusionCard:
		err = g.EvokeConfusion(ev)
	case ObstructionCard:
		err = g.EvokeObstruction(ev)
	case FireCard:
		err = g.EvokeFire(ev)
	case LignificationCard:
		err = g.EvokeLignification(ev)
	}
	if err != nil {
		return err
	}
	g.Stats.CardsUsed++ // TODO
	// TODO: animation
	g.Player.MP -= c.MPCost()
	g.FunAction()
	g.StoryPrintf("You evoked your %s.", c)
	g.DrawCard(n)
	ev.Renew(g, 5)
	return nil
}

func (c card) String() (desc string) {
	switch c {
	case BlinkCard:
		desc = "card of blinking"
	case TeleportCard:
		desc = "card of teleportation"
	case DigCard:
		desc = "card of digging"
	case TeleportOtherCard:
		desc = "card of teleport other"
	case HealWoundsCard:
		desc = "card of heal wounds"
	case MagicCard:
		desc = "card of refill magic"
	case DescentCard:
		desc = "card of descent"
	case SwiftnessCard:
		desc = "card of swiftness"
	case SwappingCard:
		desc = "card of swapping"
	case ShadowsCard:
		desc = "card of shadows"
	case FogCard:
		desc = "card of fog"
	case MagicMappingCard:
		desc = "card of magic mapping"
	case SensingCard:
		desc = "card of sensing"
	case WallsCard:
		desc = "card of walls"
	case SlowingCard:
		desc = "card of slowing"
	case SleepingCard:
		desc = "card of sleeping"
	case NoiseCard:
		desc = "card of noise"
	case ObstructionCard:
		desc = "card of obstruction"
	case ConfusionCard:
		desc = "card of confusion"
	case FireCard:
		desc = "card of fire"
	case LignificationCard:
		desc = "card of lignification"
	}
	return desc
}

func (c card) Desc(g *game) (desc string) {
	// TODO
	switch c {
	case BlinkCard:
		desc = "makes you blink away within your line of sight. The rod is more susceptible to send you to the cells thar are most far from you."
	case TeleportCard:
		desc = "makes you teleport far away."
	case DigCard:
		desc = "makes you dig walls by walking into them like an earth dragon."
	case TeleportOtherCard:
		desc = "teleports up to two random monsters in sight."
	case HealWoundsCard:
		desc = "heals you a good deal."
	case MagicCard:
		desc = "replenishes your magical reserves."
	case DescentCard:
		desc = "makes you go deeper in the Underground."
	case SwiftnessCard:
		desc = "makes you move faster and better at avoiding blows for a short time." // XXX
	case SwappingCard:
		desc = "makes you swap positions with the farthest monster in sight. If there is more than one at the same distance, it will be chosen randomly."
	case ShadowsCard:
		desc = "reduces your line of sight range to 1. Because monsters only can see you if you see them, this makes it easier to get out of sight of monsters so that they eventually stop chasing you."
	case FogCard:
		desc = ""
	case MagicMappingCard:
		desc = "shows you the map layout and item locations."
	case SensingCard:
		desc = "shows you the current position of monsters in the map."
	case WallsCard:
		desc = "replaces free cells around you with temporary walls."
	case SlowingCard:
		desc = "induces slow movement and attack for monsters in sight."
	case SleepingCard:
		desc = "induces deep sleeping and exhaustion for up to two random monsters in sight."
	case NoiseCard:
		desc = "produces a noisy bang, attracting monsters in a medium-sized area."
	case ObstructionCard:
		desc = "creates temporal walls between you and up to 3 monsters."
	case ConfusionCard:
		desc = "confuses monsters in sight, leaving them unable to attack you."
	case FireCard:
		desc = "produces a small magical fire that will extend to neighbour flammable terrain. The smoke it generates will induce sleep in monsters. As a gawalt monkey, you resist sleepiness, but you will still feel slowed."
	case LignificationCard:
		desc = "lignifies up to 2 monsters in view, so that it cannot move. The monster can still fight."
	}
	return fmt.Sprintf("The %s %s", c, desc)
}

func (c card) MPCost() int {
	if c == MagicCard {
		return 0
	}
	return 1
}

func (g *game) EvokeBlink(ev event) error {
	if g.Player.HasStatus(StatusLignification) {
		return errors.New("You cannot blink while lignified.")
	}
	g.Blink(ev)
	return nil
}

func (g *game) Blink(ev event) {
	if g.Player.HasStatus(StatusLignification) {
		return
	}
	npos := g.BlinkPos()
	if !npos.valid() {
		// should not happen
		g.Print("You could not blink.")
		return
	}
	opos := g.Player.Pos
	g.Print("You blink away.")
	g.ui.TeleportAnimation(opos, npos, true)
	g.PlacePlayerAt(npos)
}

func (g *game) BlinkPos() position {
	losPos := []position{}
	for pos, b := range g.Player.LOS {
		// TODO: skip if not seen?
		if !b {
			continue
		}
		if !g.Dungeon.Cell(pos).IsFree() {
			continue
		}
		mons := g.MonsterAt(pos)
		if mons.Exists() {
			continue
		}
		losPos = append(losPos, pos)
	}
	if len(losPos) == 0 {
		return InvalidPos
	}
	npos := losPos[RandInt(len(losPos))]
	for i := 0; i < 4; i++ {
		pos := losPos[RandInt(len(losPos))]
		if npos.Distance(g.Player.Pos) < pos.Distance(g.Player.Pos) {
			npos = pos
		}
	}
	return npos
}

func (g *game) EvokeTeleport(ev event) error {
	if g.Player.HasStatus(StatusLignification) {
		return errors.New("You cannot teleport while lignified.")
	}
	g.Teleportation(ev)
	g.Print("You teleported away.")
	return nil
}

func (g *game) EvokeDig(ev event) error {
	g.Player.Statuses[StatusDig] = 1
	end := ev.Rank() + DurationDigging
	g.PushEvent(&simpleEvent{ERank: end, EAction: DigEnd})
	g.Player.Expire[StatusDig] = end
	g.Print("You feel like an earth dragon.")
	return nil
}

func (g *game) MonstersInLOS() []*monster {
	ms := []*monster{}
	for _, mons := range g.Monsters {
		if mons.Exists() && g.Player.Sees(mons.Pos) {
			ms = append(ms, mons)
		}
	}
	// shuffle before, because the order could be unnaturally predicted
	for i := 0; i < len(ms); i++ {
		j := i + RandInt(len(ms)-i)
		ms[i], ms[j] = ms[j], ms[i]
	}
	return ms
}

func (g *game) EvokeTeleportOther(ev event) error {
	ms := g.MonstersInLOS()
	if len(ms) == 0 {
		return errors.New("There are no monsters in view.")
	}
	max := 2
	if max > len(ms) {
		max = len(ms)
	}
	for i := 0; i < max; i++ {
		ms[i].TeleportAway(g)
	}

	return nil
}

func (g *game) EvokeHealWounds(ev event) error {
	g.Player.HP = g.Player.HPMax()
	g.Print("Your feel healthy again.")
	return nil
}

func (g *game) EvokeRefillMagic(ev event) error {
	g.Player.MP = g.Player.MPMax()
	g.Print("Your magic forces return.")
	return nil
}

func (g *game) EvokeDescent(ev event) error {
	if g.Depth >= MaxDepth {
		return errors.New("You cannot descend any deeper!")
	}
	g.Printf("You fall through the ground.")
	g.LevelStats()
	g.StoryPrint("You descended deeper into the dungeon.")
	g.Depth++
	g.DepthPlayerTurn = 0
	g.InitLevel()
	g.Save()
	return nil
}

func (g *game) EvokeSwiftness(ev event) error {
	g.Player.Statuses[StatusSwift]++
	end := ev.Rank() + DurationSwiftness
	g.PushEvent(&simpleEvent{ERank: end, EAction: HasteEnd})
	g.Player.Expire[StatusSwift] = end
	g.Printf("You feel speedy and agile.")
	// XXX do something with agile?
	return nil
}

func (g *game) EvokeSwapping(ev event) error {
	if g.Player.HasStatus(StatusLignification) {
		return errors.New("You cannot use this rod while lignified.")
	}
	ms := g.MonstersInLOS()
	var mons *monster
	best := 0
	for _, m := range ms {
		if m.Status(MonsLignified) {
			continue
		}
		if m.Pos.Distance(g.Player.Pos) > best {
			best = m.Pos.Distance(g.Player.Pos)
			mons = m
		}
	}
	if !mons.Exists() {
		return errors.New("No monsters suitable for swapping in view.")
	}
	g.SwapWithMonster(mons)
	return nil
}

func (g *game) SwapWithMonster(mons *monster) {
	ompos := mons.Pos
	g.Printf("You swap positions with the %s.", mons.Kind)
	g.ui.SwappingAnimation(mons.Pos, g.Player.Pos)
	mons.MoveTo(g, g.Player.Pos)
	g.PlacePlayerAt(ompos)
	mons.MakeAware(g)
}

func (g *game) EvokeShadows(ev event) error {
	if g.Player.HasStatus(StatusShadows) {
		return errors.New("You are already surrounded by shadows.")
	}
	g.Player.Statuses[StatusShadows] = 1
	end := ev.Rank() + DurationShadows
	g.PushEvent(&simpleEvent{ERank: end, EAction: ShadowsEnd})
	g.Player.Expire[StatusShadows] = end
	g.Printf("You feel surrounded by shadows.")
	g.ComputeLOS()
	return nil
}

type cloud int

const (
	CloudFog cloud = iota
	CloudFire
	CloudNight
)

func (g *game) EvokeFog(ev event) error {
	g.Fog(g.Player.Pos, 3, ev)
	g.Print("You are surrounded by a dense fog.")
	return nil
}

func (g *game) Fog(at position, radius int, ev event) {
	dij := &normalPath{game: g}
	nm := Dijkstra(dij, []position{at}, radius)
	for pos := range nm {
		_, ok := g.Clouds[pos]
		if !ok {
			g.Clouds[pos] = CloudFog
			g.PushEvent(&cloudEvent{ERank: ev.Rank() + DurationFog + RandInt(DurationFog/2), EAction: CloudEnd, Pos: pos})
		}
	}
	g.ComputeLOS()
}

func (g *game) EvokeMagicMapping(ev event) error {
	dp := &dungeonPath{dungeon: g.Dungeon}
	g.AutoExploreDijkstra(dp, []int{g.Player.Pos.idx()})
	cdists := make(map[int][]int)
	for i, dist := range DijkstraMapCache {
		cdists[dist] = append(cdists[dist], i)
	}
	var dists []int
	for dist, _ := range cdists {
		dists = append(dists, dist)
	}
	sort.Ints(dists)
	g.ui.DrawDungeonView(NormalMode)
	for _, d := range dists {
		draw := false
		for _, i := range cdists[d] {
			pos := idxtopos(i)
			c := g.Dungeon.Cell(pos)
			if (c.IsFree() || g.Dungeon.HasFreeNeighbor(pos)) && !c.Explored {
				g.Dungeon.SetExplored(pos)
				draw = true
			}
		}
		if draw {
			g.ui.MagicMappingAnimation(cdists[d])
		}
	}
	g.Printf("You feel aware of your surroundings..")
	return nil
}

func (g *game) EvokeSensing(ev event) error {
	for _, mons := range g.Monsters {
		if mons.Exists() && !g.Player.Sees(mons.Pos) {
			mons.UpdateKnowledge(g, mons.Pos)
		}
	}
	g.Printf("You briefly sense monsters around.")
	return nil
}

func (g *game) EvokeWalls(ev event) error {
	neighbors := g.Dungeon.FreeNeighbors(g.Player.Pos)
	for _, pos := range neighbors {
		mons := g.MonsterAt(pos)
		if mons.Exists() {
			continue
		}
		g.CreateTemporalWallAt(pos, ev)
	}
	g.Printf("You feel surrounded by temporary walls.")
	g.ComputeLOS()
	return nil
}

func (g *game) EvokeSlowing(ev event) error {
	g.Print("Whoosh! A slowing luminous wave emerges.")
	// TODO: animation
	for _, mons := range g.Monsters {
		if !mons.Exists() || !g.Player.Sees(mons.Pos) {
			continue
		}
		mons.Statuses[MonsSlow]++
		g.PushEvent(&monsterEvent{ERank: g.Ev.Rank() + DurationSlow, NMons: mons.Index, EAction: MonsSlowEnd})
	}

	ev.Renew(g, DurationThrowItem)
	return nil
}

func (g *game) EvokeSleeping(ev event) error {
	ms := g.MonstersInLOS()
	if len(ms) == 0 {
		return errors.New("There are no monsters in view.")
	}
	max := 2
	if max > len(ms) {
		max = len(ms)
	}
	g.Print("Two beams of sleeping emerge.")
	// TODO: animation
	// XXX: maybe use noise distance instead of LOS?
	for i := 0; i < max; i++ {
		mons := ms[i]
		if mons.State != Resting {
			g.Printf("%s falls asleep.", mons.Kind.Definite(true))
		}
		mons.State = Resting
		mons.Dir = NoDir
		mons.ExhaustTime(g, 40+RandInt(10))
	}
	return nil
}

func (g *game) EvokeLignification(ev event) error {
	ms := g.MonstersInLOS()
	if len(ms) == 0 {
		return errors.New("There are no monsters in view.")
	}
	max := 2
	if max > len(ms) {
		max = len(ms)
	}
	g.Print("Two beams of lignification emerge.")
	// TODO: animation
	for i := 0; i < max; i++ {
		mons := ms[i]
		if mons.Status(MonsLignified) {
			continue
		}
		mons.EnterLignification(g, ev)
	}
	return nil
}

func (g *game) EvokeNoise(ev event) error {
	g.MakeNoise(CardBangNoise, g.Player.Pos)
	g.Print("Baaang!!! You better get out of here.")
	return nil
}

func (g *game) EvokeConfusion(ev event) error {
	g.Print("Whoosh! A slowing luminous wave emerges.")
	// TODO: animation
	for _, mons := range g.Monsters {
		if !mons.Exists() || !g.Player.Sees(mons.Pos) {
			continue
		}
		mons.Statuses[MonsConfused]++
		g.PushEvent(&monsterEvent{ERank: g.Ev.Rank() + DurationConfusion, NMons: mons.Index, EAction: MonsSlowEnd})
	}

	ev.Renew(g, DurationThrowItem)
	return nil
}

func (g *game) EvokeFire(ev event) error {
	if !g.Dungeon.Cell(g.Player.Pos).Flammable() {
		return errors.New("You are not standing on flammable terrain.")
	}
	g.Burn(g.Player.Pos, ev)
	g.Print("A small fire appears.")
	return nil
}

func (g *game) EvokeObstruction(ev event) error {
	ms := g.MonstersInLOS()
	if len(ms) == 0 {
		return errors.New("There are no monsters in view.")
	}
	max := 3
	if max > len(ms) {
		max = len(ms)
	}
	for i := 0; i < max; i++ {
		ray := g.Ray(ms[i].Pos)
		for _, pos := range ray[1:] {
			if pos == g.Player.Pos {
				break
			}
			mons := g.MonsterAt(pos)
			if mons.Exists() {
				continue
			}
			g.TemporalWallAt(pos, ev)
			break
		}
	}
	g.Print("Magical walls emerged.")
	return nil
}

func (g *game) TemporalWallAt(pos position, ev event) {
	if g.Dungeon.Cell(pos).T == WallCell {
		return
	}
	if !g.Player.Sees(pos) {
		g.TerrainKnowledge[pos] = g.Dungeon.Cell(pos).T
	}
	g.CreateTemporalWallAt(pos, ev)
	g.ComputeLOS()
}

func (g *game) CreateTemporalWallAt(pos position, ev event) {
	t := g.Dungeon.Cell(pos).T
	g.Dungeon.SetCell(pos, WallCell)
	delete(g.Clouds, pos)
	g.TemporalWalls[pos] = t
	g.PushEvent(&cloudEvent{ERank: ev.Rank() + DurationTemporalWall + RandInt(DurationTemporalWall/2), Pos: pos, EAction: ObstructionEnd})
}

////////////////////////////////////////////////////////////////

//func (g *game) QuaffBerserk(ev event) error {
//if g.Player.HasStatus(StatusExhausted) {
//return errors.New("You are too exhausted to berserk.")
//}
//if g.Player.HasStatus(StatusBerserk) {
//return errors.New("You are already berserk.")
//}
//g.Player.Statuses[StatusBerserk] = 1
//end := ev.Rank() + DurationBerserk
//g.PushEvent(&simpleEvent{ERank: end, EAction: BerserkEnd})
//g.Player.Expire[StatusBerserk] = end
//g.Printf("You quaff the %s. You feel a sudden urge to kill things.", BerserkPotion)
//g.Player.HPbonus += 2
//return nil
//}

//func (g *game) QuaffSwapPotion(ev event) error {
//if g.Player.HasStatus(StatusLignification) {
//return errors.New("You cannot drink this potion while lignified.")
//}
//g.Player.Statuses[StatusSwap] = 1
//end := ev.Rank() + DurationSwap
//g.PushEvent(&simpleEvent{ERank: end, EAction: SwapEnd})
//g.Player.Expire[StatusSwap] = end
////g.Printf("You quaff the %s. You feel light-footed.", SwapPotion)
//return nil
//}

//func (g *game) QuaffLignification(ev event) error {
//if g.Player.HasStatus(StatusLignification) {
//return errors.New("You are already lignified.")
//}
//g.EnterLignification(ev)
//g.Printf("You quaff the %s. You feel rooted to the ground.", LignificationPotion)
//return nil
//}

//func (g *game) QuaffCBlinkPotion(ev event) error {
//if g.Player.HasStatus(StatusLignification) {
//return errors.New("You cannot blink while lignified.")
//}
//if err := g.ui.ChooseTarget(&chooser{free: true}); err != nil {
//return err
//}
////g.Printf("You quaff the %s. You blink.", CBlinkPotion)
//g.PlacePlayerAt(g.Player.Target)
//return nil
//}

//func (g *game) ExplosionAt(ev event, pos position) {
//g.Burn(pos, ev)
//mons := g.MonsterAt(pos)
//if mons.Exists() {
//mons.HP /= 2
//if mons.HP <= 0 {
//g.HandleKill(mons, ev)
//if g.Player.Sees(mons.Pos) {
//g.Printf("%s dies.", mons.Kind.Definite(true))
//}
//}
//g.MakeNoise(ExplosionHitNoise, mons.Pos)
//g.HandleStone(mons)
//mons.MakeHuntIfHurt(g)
//} else if c := g.Dungeon.Cell(pos); !c.IsFree() && RandInt(2) == 0 {
//g.Dungeon.SetCell(pos, GroundCell)
//g.Stats.Digs++
//if !g.Player.Sees(pos) {
//g.TerrainKnowledge[pos] = c.T
//} else {
//g.ui.WallExplosionAnimation(pos)
//}
//g.MakeNoise(WallNoise, pos)
//g.Fog(pos, 1, ev)
//}
//}

//func (g *game) ThrowExplosiveMagara(ev event) error {
//if err := g.ui.ChooseTarget(&chooser{area: true, minDist: true, flammable: true, wall: true}); err != nil {
//return err
//}
//neighbors := g.Player.Target.ValidNeighbors()
//g.Printf("You throw the explosive magara... %s", g.ExplosionSound())
//g.MakeNoise(ExplosionNoise, g.Player.Target)
//g.ui.ProjectileTrajectoryAnimation(g.Ray(g.Player.Target), ColorFgPlayer)
//g.ui.ExplosionAnimation(FireExplosion, g.Player.Target)
//for _, pos := range append(neighbors, g.Player.Target) {
//g.ExplosionAt(ev, pos)
//}

//ev.Renew(g, DurationThrowItem)
//return nil
//}

//func (g *game) ThrowTeleportMagara(ev event) error {
//if err := g.ui.ChooseTarget(&chooser{area: true, minDist: true}); err != nil {
//return err
//}
//neighbors := g.Player.Target.ValidNeighbors()
//g.Print("You throw the teleport magara.")
//g.ui.ProjectileTrajectoryAnimation(g.Ray(g.Player.Target), ColorFgPlayer)
//for _, pos := range append(neighbors, g.Player.Target) {
//mons := g.MonsterAt(pos)
//if mons.Exists() {
//mons.TeleportAway(g)
//}
//}

//ev.Renew(g, DurationThrowItem)
//return nil
//}

//func (g *game) ThrowNightMagara(ev event) error {
//if err := g.ui.ChooseTarget(&chooser{needsFreeWay: true}); err != nil {
//return err
//}
//g.Print("You throw the night magara… Clouds come out of it.")
//g.ui.ProjectileTrajectoryAnimation(g.Ray(g.Player.Target), ColorFgSleepingMonster)
//g.NightFog(g.Player.Target, 2, ev)

//ev.Renew(g, DurationThrowItem)
//return nil
//}

//func (g *game) EvokeRodFireBolt(ev event) error {
//if err := g.ui.ChooseTarget(&chooser{flammable: true}); err != nil {
//return err
//}
//ray := g.Ray(g.Player.Target)
//g.MakeNoise(MagicCastNoise, g.Player.Pos)
//g.Print("Whoosh! A fire bolt emerges straight out of the rod.")
//g.ui.FireBoltAnimation(ray)
//for _, pos := range ray {
//g.Burn(pos, ev)
//mons := g.MonsterAt(pos)
//if !mons.Exists() {
//continue
//}
//dmg := 1
//mons.HP -= dmg
//if mons.HP <= 0 {
//g.Printf("%s is killed by the bolt.", mons.Kind.Indefinite(true))
//g.HandleKill(mons, ev)
//}
//g.MakeNoise(MagicHitNoise, mons.Pos)
//g.HandleStone(mons)
//mons.MakeHuntIfHurt(g)
//}
//return nil
//}

//func (g *game) EvokeRodFireball(ev event) error {
//if err := g.ui.ChooseTarget(&chooser{area: true, minDist: true, flammable: true}); err != nil {
//return err
//}
//neighbors := g.Dungeon.FreeNeighbors(g.Player.Target)
//g.MakeNoise(MagicExplosionNoise, g.Player.Target)
//g.Printf("A fireball emerges straight out of the rod... %s", g.ExplosionSound())
//g.ui.ProjectileTrajectoryAnimation(g.Ray(g.Player.Target), ColorFgExplosionStart)
//g.ui.ExplosionAnimation(FireExplosion, g.Player.Target)
//for _, pos := range append(neighbors, g.Player.Target) {
//g.Burn(pos, ev)
//mons := g.MonsterAt(pos)
//if mons == nil {
//continue
//}
//dmg := 1 + RandInt(2)
//mons.HP -= dmg
//if mons.HP <= 0 {
//g.Printf("%s is killed by the fireball.", mons.Kind.Indefinite(true))
//g.HandleKill(mons, ev)
//}
//g.MakeNoise(MagicHitNoise, mons.Pos)
//g.HandleStone(mons)
//mons.MakeHuntIfHurt(g)
//}
//return nil
//}

//func (g *game) EvokeRodLightning(ev event) error {
//d := g.Dungeon
//conn := map[position]bool{}
//nb := make([]position, 0, 8)
//nb = g.Player.Pos.Neighbors(nb, func(npos position) bool {
//return npos.valid() && d.Cell(npos).T != WallCell
//})
//stack := []position{}
//g.MakeNoise(MagicCastNoise, g.Player.Pos)
//g.Print("Whoosh! Lightning emerges straight out of the rod.")
//for _, pos := range nb {
//mons := g.MonsterAt(pos)
//if !mons.Exists() {
//continue
//}
//stack = append(stack, pos)
//conn[pos] = true
//}
//if len(stack) == 0 {
//return errors.New("There are no adjacent monsters.")
//}
//var pos position
//targets := []position{}
//for len(stack) > 0 {
//pos = stack[len(stack)-1]
//stack = stack[:len(stack)-1]
//g.Burn(pos, ev)
//mons := g.MonsterAt(pos)
//if !mons.Exists() {
//continue
//}
//targets = append(targets, pos)
//dmg := 1
//mons.HP -= dmg
//if mons.HP <= 0 {
//g.Printf("%s is killed by lightning.", mons.Kind.Indefinite(true))
//g.HandleKill(mons, ev)
//}
//g.MakeNoise(MagicHitNoise, mons.Pos)
//g.HandleStone(mons)
//mons.MakeHuntIfHurt(g)
//nb = pos.Neighbors(nb, func(npos position) bool {
//return npos.valid() && d.Cell(npos).T != WallCell
//})
//for _, npos := range nb {
//if !conn[npos] {
//conn[npos] = true
//stack = append(stack, npos)
//}
//}
//}
//g.ui.LightningHitAnimation(targets)

//return nil
//}

//func (g *game) EvokeRodDigging(ev event) error {
//if err := g.ui.ChooseTarget(&wallChooser{}); err != nil {
//return err
//}
//pos := g.Player.Target
//for i := 0; i < 3; i++ {
//g.Dungeon.SetCell(pos, GroundCell)
//g.Stats.Digs++
//g.MakeNoise(WallNoise, pos)
//g.Fog(pos, 1, ev)
//pos = pos.To(pos.Dir(g.Player.Pos))
//if !g.Player.Sees(pos) {
//g.TerrainKnowledge[pos] = WallCell
//}
//if !pos.valid() || g.Dungeon.Cell(pos).T != WallCell {
//break
//}
//}
//g.Print("You see the wall disintegrate with a crash.")
//g.ComputeLOS()
//g.MakeMonstersAware()
//return nil
//}

//func (g *game) EvokeRodShatter(ev event) error {
//if err := g.ui.ChooseTarget(&wallChooser{minDist: true}); err != nil {
//return err
//}
//neighbors := g.Dungeon.FreeNeighbors(g.Player.Target)
//g.Dungeon.SetCell(g.Player.Target, GroundCell)
//g.Stats.Digs++
//g.ComputeLOS()
//g.MakeMonstersAware()
//g.MakeNoise(WallNoise, g.Player.Target)
//g.Printf("%s The wall disappeared.", g.CrackSound())
//g.ui.ProjectileTrajectoryAnimation(g.Ray(g.Player.Target), ColorFgExplosionWallStart)
//g.ui.ExplosionAnimation(WallExplosion, g.Player.Target)
//g.Fog(g.Player.Target, 2, ev)
//for _, pos := range neighbors {
//mons := g.MonsterAt(pos)
//if !mons.Exists() {
//continue
//}
//dmg := 2
//mons.HP -= dmg
//if mons.HP <= 0 {
//g.Printf("%s is killed by the explosion.", mons.Kind.Indefinite(true))
//g.HandleKill(mons, ev)
//}
//g.MakeNoise(ExplosionHitNoise, mons.Pos)
//g.HandleStone(mons)
//mons.MakeHuntIfHurt(g)
//}
//return nil
//}

//func (g *game) EvokeRodLignification(ev event) error {
//if err := g.ui.ChooseTarget(&chooser{}); err != nil {
//return err
//}
//mons := g.MonsterAt(g.Player.Target)
//// mons not nil (check done in targeter)
//if mons.Status(MonsLignified) {
//return errors.New("You cannot target a lignified monster.")
//}
//mons.EnterLignification(g, ev)
//return nil
//}

//func (g *game) EvokeRodHope(ev event) error {
//if err := g.ui.ChooseTarget(&chooser{needsFreeWay: true}); err != nil {
//return err
//}
//g.MakeNoise(MagicCastNoise, g.Player.Pos)
//g.ui.ProjectileTrajectoryAnimation(g.Ray(g.Player.Target), ColorFgExplosionStart)
//mons := g.MonsterAt(g.Player.Target)
//// mons not nil (check done in the targeter)
//dmg := DefaultHealth - g.Player.HP + 1
//if dmg <= 0 {
//dmg = 1
//}
//mons.HP -= dmg
//g.Burn(g.Player.Target, ev)
//g.ui.HitAnimation(g.Player.Target, true)
//g.Printf("An energy channel hits %s (%d dmg).", mons.Kind.Definite(false), dmg)
//if mons.HP <= 0 {
//g.Printf("%s dies.", mons.Kind.Indefinite(true))
//g.HandleKill(mons, ev)
//}
//return nil
//}