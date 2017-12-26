package model

import (
	"adventure/advserver/gamedata"
	"adventure/common/structs"
	"adventure/common/util"
)

type FightSim struct {
	Left          *structs.FightTeam
	Right         *structs.FightTeam
	MaxRoundCount int32
}

func NewFightSim(left, right *structs.FightTeam) *FightSim {
	sim := &FightSim{left, right, 99}
	return sim
}

func getSkills(skillIDs []int32) []*structs.SpellTemplate {
	skills := []*structs.SpellTemplate{}
	for _, v := range skillIDs {
		skill, ok := gamedata.AllTemplates.SpellTemplates[v]
		if ok {
			skills = append(skills, &skill)
		}
	}
	return skills
}

func (f *FightSim) Fight() *structs.FightResult {

	isLeftRound := f.Left.XianGong >= f.Right.XianGong
	// 左边
	leftHP := f.Left.DefaultHP                                              // 血量
	leftMissRate := float32(f.Left.ShanBi) / float32(f.Left.ShanBi+100)     // 闪避率
	leftResistRate := 1 - float32(f.Left.FangYu)/float32(f.Left.FangYu+100) // 免伤率
	leftDeepRate := float32(f.Left.WangZhe)*0.01 + 1                        // 易伤率
	leftSkills := getSkills(f.Left.SpellIDs)

	// 右边
	rightHP := f.Right.DefaultHP                                               // 血量
	rightMissRate := float32(f.Right.ShanBi) / float32(f.Right.ShanBi+100)     // 闪避率
	rightResistRate := 1 - float32(f.Right.FangYu)/float32(f.Right.FangYu+100) // 免伤率
	rightDeepRate := float32(f.Right.WangZhe)*0.01 + 1                         // 易伤率
	rightSkills := getSkills(f.Right.SpellIDs)

	userSkillCount := make(map[int32]int32)
	rounds := []structs.FightRound{}
	roundIndex := int32(0)
	for leftHP > 0 && rightHP > 0 && roundIndex < f.MaxRoundCount {
		canUseSkills := []*structs.SpellTemplate{}

		curSkills := rightSkills
		if isLeftRound {
			curSkills = leftSkills
		}
		if len(curSkills) > 0 {
			for i := 0; i < len(curSkills); i++ {
				s2 := curSkills[i]

				// skill index
				randIndex := 10000
				if isLeftRound {
					randIndex = 200
				}
				skillIndex := int32(i + randIndex)
				bContain := false
				for _, v := range userSkillCount {
					if v == skillIndex {
						bContain = true
					}
				}
				if !bContain {
					userSkillCount[skillIndex] = 0
				}
				if skillCanUse(isLeftRound, roundIndex, f.Left.DefaultHP, f.Right.DefaultHP, leftHP, rightHP, userSkillCount[skillIndex], s2) {
					canUseSkills = append(canUseSkills, s2)
				}
			}
		}

		skill := new(structs.SpellTemplate)
		randomValue := int32(util.RandNum(0, 100))
		curValue := int32(0)
		for _, s3 := range canUseSkills {
			curValue += s3.Rate
			if randomValue < curValue {
				skill = s3
				userSkillCount[s3.ID]++
			}
		}

		ignoreDefence := false
		isMiss := false
		damage := int32(0)
		recover := int32(0)
		randValue := float32(0)

		missRate := leftMissRate
		if isLeftRound {
			missRate = rightMissRate
		}
		if skill == nil {

			ignoreDefence = false

			damage = f.Right.DefaultHP
			randValue = float32(util.RandNum(gamedata.FightFloatValueMin, gamedata.FightFloatValueMax))
			if isLeftRound {
				damage = f.Left.DefaultHP
			}
			damage = int32(float32(damage) * 0.01 * float32(randValue))

			recover = 0

			if float32(util.RandNum(0, 100))/100 < missRate {
				isMiss = true
			}

		} else {
			if !skill.IgnoreDodge && float32(util.RandNum(0, 100)) < missRate {
				isMiss = true
			}

			ignoreDefence = skill.IgnoreDefence

			randValue = float32(util.RandNum(gamedata.FightSkillFloatValueMin, gamedata.FightSkillFloatValueMax))
			if skill.AttackType == structs.Hurt {
				damage = computeValue(roundIndex, f.Left.DefaultHP, f.Right.DefaultHP, leftHP, rightHP, 0, isLeftRound, true, skill, randValue)
			} else if skill.AttackType == structs.Recover {
				damage = 0
				recover = computeValue(roundIndex, f.Left.DefaultHP, f.Right.DefaultHP, leftHP, rightHP, 0, isLeftRound, false, skill, randValue)

			} else if skill.AttackType == structs.HurtAndRecover {
				damage = computeValue(roundIndex, f.Left.DefaultHP, f.Right.DefaultHP, leftHP, rightHP, 0, isLeftRound, true, skill, randValue)
				if skill.CalculateRecoverType != structs.ThisTimeAttack {
					recover = computeValue(roundIndex, f.Left.DefaultHP, f.Right.DefaultHP, leftHP, rightHP, damage, isLeftRound, false, skill, randValue)
				}
			}
		}

		damage2 := int32(0) // 最终伤害
		if !isMiss {
			resistRate := leftResistRate
			if isLeftRound {
				resistRate = rightResistRate
			}
			if ignoreDefence {
				resistRate = 1
			}
			deepRate := rightDeepRate
			if isLeftRound {
				deepRate = leftDeepRate
			}
			damage2 = int32(float32(damage) * resistRate * deepRate)

			if skill != nil && skill.AttackType == structs.HurtAndRecover && skill.CalculateRecoverType == structs.ThisTimeAttack {
				recover = int32(damage2 / 100 * skill.RecoverEffectValue)
			}
		}

		if isLeftRound {
			leftHP = leftHP + recover
			rightHP = rightHP - damage2
		} else {
			rightHP = rightHP + recover
			leftHP = leftHP - damage2
		}

		skillID := int32(-1)
		if skill != nil {
			skillID = skill.ID
		}

		// 添加回合数据
		rounds = append(rounds, structs.FightRound{skillID, leftHP, rightHP})
		// 切换队伍
		isLeftRound = !isLeftRound
		// 添加回合数
		roundIndex++
	}

	fightRet := &structs.FightResult{
		LeftTeam:  f.Left,
		RightTeam: f.Right,
		Rounds:    rounds,
		IsLeftWin: leftHP > rightHP,
	}

	return fightRet
}

func skillCanUse(isLeftRound bool, round, leftDefaultHP, rightDefaultHP, leftHP, rightHP, userCount int32, skill *structs.SpellTemplate) bool {
	result := true
	if skill.RoundTriggerType == structs.RoundMax {
		result = round+1 > skill.RoundValue
	} else if skill.RoundTriggerType == structs.RoundMin {
		result = round+1 < skill.RoundValue
	}

	if !result {
		return false
	}

	leftCapacity := float32(leftHP*100) / float32(leftDefaultHP)
	rightCapacity := float32(rightHP*100) / float32(rightDefaultHP)

	switch skill.FightingCapacityTriggerType {
	case structs.EnemyHPmin:
		capacity := leftCapacity
		if isLeftRound {
			capacity = rightCapacity
		}
		result = int32(capacity) < skill.FightingCapacity
	case structs.SelfHPmin:
		capacity := rightCapacity
		if isLeftRound {
			capacity = leftCapacity
		}
		result = int32(capacity) < skill.FightingCapacity
	case structs.EnemyHPmax:
		capacity := leftCapacity
		if isLeftRound {
			capacity = rightCapacity
		}
		result = int32(capacity) > skill.FightingCapacity
	case structs.SelfHPmax:
		capacity := rightCapacity
		if isLeftRound {
			capacity = leftCapacity
		}
		result = int32(capacity) > skill.FightingCapacity
	}
	if !result {
		return false
	}

	if skill.CountType == structs.One {
		result = userCount < 1
	}

	return result
}

func computeValue(round, leftDefaultHP, rightDefaultHP, leftHP, rightHP, theDamage int32, isLeftRound, isDamage bool, skill *structs.SpellTemplate, randValue float32) int32 {
	value := int32(0)

	skillType := skill.CalculateRecoverType
	if isDamage {
		skillType = skill.CalculateHurtType
	}
	switch skillType {
	case structs.EnemyLostFight:
		value = leftDefaultHP - leftHP
		if isLeftRound {
			value = rightDefaultHP - rightHP
		}
		if value < 0 {
			value = 1
		}
	case structs.EnemyTotalFight:
		value = leftDefaultHP
		if isLeftRound {
			value = rightDefaultHP
		}
	case structs.EnumyLeftFight:
		value = leftHP
		if isLeftRound {
			value = rightHP
		}
	case structs.SelfLeftFight:
		value = rightHP
		if isLeftRound {
			value = leftHP
		}
	case structs.SelfLostFight:
		value = rightDefaultHP - rightHP
		if isLeftRound {
			value = leftDefaultHP - leftHP
		}
	case structs.SelfTotalFight:
		value = rightDefaultHP
		if isLeftRound {
			value = leftDefaultHP
		}
	case structs.ThisTimeAttack:
		value = theDamage
	case structs.NormalAttack:
		value = rightDefaultHP
		if isLeftRound {
			value = leftDefaultHP
		}
		value = int32(float32(value) / 200 * (gamedata.FightFloatValueMin + gamedata.FightFloatValueMax))
	}

	damage := skill.RecoverEffectValue
	if isDamage {
		damage = skill.HurtEffectValue
	}

	result := int32(float32(value) * 0.01 * float32(damage) * randValue * 0.01)

	if result <= 1 {
		result = 1
	}

	return result
}
